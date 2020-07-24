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
	Bots       map[uint]*Bot
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
		Bots:       make(map[uint]*Bot),
		Logger:     log,
	}
}

// About -
func (c *Client) About() (string, error) {
	res, err := c.doGet("/about", nil)
	if err != nil {
		return "", err
	}
	return gjson.Get(res, "data.version").String(), nil
}

// Auth 开始会话-认证(Authorize)
func (c *Client) Auth() (string, error) {
	data := map[string]string{"authKey": c.AuthKey}
	res, err := c.doPost("/auth", data)
	if err != nil {
		return "", err
	}
	c.Logger.Infoln("Authed")
	return gjson.Get(res, "session").String(), nil
}

// Verify 校验Session
func (c *Client) Verify(qq uint, sessionKey string) (*Bot, error) {
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

// Release 释放Session
func (c *Client) Release(qq uint) error {
	data := map[string]interface{}{"sessionKey": c.Bots[qq].SessionKey, "qq": qq}
	_, err := c.doPost("release", data)
	if err != nil {
		return err
	}
	delete(c.Bots, qq)
	c.Logger.Infoln("Released")
	return nil
}

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
	return res.String(), getErrByCode(uint(gjson.Get(res.String(), "code").Uint()))
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
	return res.String(), getErrByCode(uint(gjson.Get(res.String(), "code").Uint()))
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
	return res.String(), getErrByCode(uint(gjson.Get(res.String(), "code").Uint()))
}

func getErrByCode(code uint) error {
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
