package gomirai

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/body"
	"gopkg.in/h2non/gentleman.v2/plugins/multipart"

	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

// Client 与Mirai进行沟通
type Client struct {
	Name       string
	AuthKey    string
	HTTPClient *gentleman.Client
	Bots       map[int64]*Bot
	Logger     *logrus.Entry
}

// NewClient 新建Client
func NewClient(name, url, authKey string, logger ...*logrus.Entry) *Client {
	c := gentleman.New()
	c.URL(url)

	var log *logrus.Entry
	if len(logger) > 0 {
		log = logger[0]
	} else {
		log = logrus.New().WithFields(logrus.Fields{
			"client": name,
		})
	}
	return &Client{
		AuthKey:    authKey,
		HTTPClient: c,
		Bots:       make(map[int64]*Bot),
		Logger:     log,
	}
}

// --- API-HTTP插件相关 ---

// About 使用此方法获取插件的信息，如版本号
func (c *Client) About() (string, error) {
	res, err := c.doGet("/about", nil)
	if err != nil {
		return "", err
	}
	return gjson.Get(res, "data.version").String(), nil
}

// --- 认证相关 ---

// Auth 使用此方法验证你的身份，并返回一个会话
func (c *Client) Auth() (string, error) {
	data := map[string]string{"authKey": c.AuthKey}
	res, err := c.doPost("/auth", data)
	if err != nil {
		return "", err
	}
	c.Logger.Infoln("Authed")
	return gjson.Get(res, "session").String(), nil
}

// Verify 使用此方法校验并激活你的Session，同时将Session与一个已登录的Bot绑定
func (c *Client) Verify(qq int64, sessionKey string) (*Bot, error) {
	data := map[string]interface{}{"sessionKey": sessionKey, "qq": qq}
	_, err := c.doPost("/verify", data)
	if err != nil {
		return nil, err
	}
	c.Bots[qq] = &Bot{QQ: qq, SessionKey: sessionKey, Client: c, Logger: c.Logger.WithField("qq", qq)}
	c.Bots[qq].SetChannel(time.Second, 10)
	c.Logger.Infoln("Verified")
	return c.Bots[qq], nil
}

// Release 使用此方式释放session及其相关资源（Bot不会被释放）
// 不使用的Session应当被释放，长时间（30分钟）未使用的Session将自动释放，否则Session持续保存Bot收到的消息，将会导致内存泄露(开启websocket后将不会自动释放)
func (c *Client) Release(qq int64) error {
	data := map[string]interface{}{"sessionKey": c.Bots[qq].SessionKey, "qq": qq}
	_, err := c.doPost("release", data)
	if err != nil {
		return err
	}
	delete(c.Bots, qq)
	c.Logger.Infoln("Released")
	return nil
}

// --- internal ---

func (c *Client) doPost(path string, data interface{}) (string, error) {
	c.Logger.Debugln("POST:", path, " Data:", data)
	res, err := c.HTTPClient.Request().
		Path(path).
		Method("POST").
		Use(body.JSON(data)).
		SetHeader("Content-Type", "application/json;charset=utf-8").
		Send()
	if err != nil {
		c.Logger.Warn("POST Failed")
		return "", err
	}
	c.Logger.Trace("result StatusCode:", res.StatusCode)
	if !res.Ok {
		return res.String(), errors.New("Http: " + strconv.Itoa(res.StatusCode))
	}
	return res.String(), getErrByCode(gjson.Get(res.String(), "code").Int())
}

func (c *Client) doPostWithFormData(path string, fields map[string]interface{}) (string, error) {
	data := make(multipart.DataFields)
	files := make([]multipart.FormFile, 0)
	for key, value := range fields {
		if unbox, ok := value.(string); ok {
			data[key] = append(data[key], unbox)
		} else if unbox, ok := value.(io.Reader); ok {
			files = append(files, multipart.FormFile{Name: key, Reader: unbox})
		}
	}
	formData := multipart.FormData{Data: data, Files: files}
	c.Logger.Trace("POST:"+path+" FormData:", formData)
	res, err := c.HTTPClient.Request().
		Path(path).
		Method("POST").
		Use(multipart.Data(formData)).
		Send()
	if err != nil {
		c.Logger.Warn("POST Failed")
		return "", err
	}
	c.Logger.Debugln("result StatusCode:", res.StatusCode)
	if !res.Ok {
		return res.String(), fmt.Errorf("HTTP: %d", res.StatusCode)
	}
	return res.String(), getErrByCode(gjson.Get(res.String(), "code").Int())
}

func (c *Client) doGet(path string, params map[string]string) (string, error) {
	c.Logger.Debugln("GET:", path)
	res, err := c.HTTPClient.Request().
		Path(path).
		SetQueryParams(params).
		Method("GET").
		SetHeader("Content-Type", "application/json;charset=utf-8").
		Send()
	if err != nil {
		c.Logger.Warn("GET Failed")
		return "", err
	}
	c.Logger.Debugln("result StatusCode:", res.StatusCode)
	if !res.Ok {
		return res.String(), fmt.Errorf("HTTP: %d", res.StatusCode)
	}
	return res.String(), getErrByCode(gjson.Get(res.String(), "code").Int())
}

func getErrByCode(code int64) error {
	switch code {
	case 0:
		return nil
	case 1:
		return fmt.Errorf("错误的auth key")
	case 2:
		return fmt.Errorf("指定的Bot不存在")
	case 3:
		return fmt.Errorf("Session失效或不存在")
	case 4:
		return fmt.Errorf("Session未认证(未激活)")
	case 5:
		return fmt.Errorf("发送消息目标不存在(指定对象不存在)")
	case 6:
		return fmt.Errorf("指定文件不存在，出现于发送本地图片")
	case 10:
		return fmt.Errorf("无操作权限，指Bot没有对应操作的限权")
	case 20:
		return fmt.Errorf("Bot被禁言，指Bot当前无法向指定群发送消息")
	case 30:
		return fmt.Errorf("消息过长")
	case 400:
		return fmt.Errorf("错误的访问，如参数错误等")
	default:
		return fmt.Errorf("未知错误，Code: %d", code)
	}
}
