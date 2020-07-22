package gomirai

import (
	"fmt"
	"time"

	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/body"

	"github.com/sirupsen/logrus"
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

// Auth 开始会话-认证(Authorize)
func (c *Client) Auth() (string, error) {
	data := map[string]string{"authKey": c.AuthKey}
	res, err := c.doPost("/auth", data)
	if err != nil {
		return "", err
	}
	c.Logger.Infoln("Authed")
	return JSON.Get([]byte(res), "session").ToString(), nil
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
		SetHeader("Content-Type", "application/json;charset=utf-8").
		Use(body.JSON(data)).
		Send()
	if err != nil {
		c.Logger.Warn("POST Failed")
		return "", err
	}
	c.Logger.Debugln("result StatusCode:", res.StatusCode)
	if !res.Ok {
		return res.String(), fmt.Errorf("HTTP: %d", res.StatusCode)
		// errors.New("Http: " + strconv.Itoa(res.StatusCode))
	}
	if JSON.Get([]byte(res.String()), "code").ToInt() != 0 {
		return res.String(), getErrByCode(JSON.Get([]byte(res.String()), "code").ToUint())
	}
	return res.String(), nil
}

func (c *Client) doGet(path string, params map[string]string) (string, error) {
	c.Logger.Debugln("GET:", path)
	res, err := c.HTTPClient.Request().
		Path(path).
		SetQueryParams(params).
		SetHeader("Content-Type", "application/json;charset=utf-8").
		Method("GET").
		Send()
	if err != nil {
		c.Logger.Warn("GET Failed")
		return "", err
	}
	c.Logger.Debugln("result StatusCode:", res.StatusCode)
	if !res.Ok {
		return res.String(), fmt.Errorf("HTTP: %d", res.StatusCode)
		// errors.New("Http: " + strconv.Itoa(res.StatusCode))
	}
	if JSON.Get([]byte(res.String()), "code").ToInt() != 0 {
		return res.String(), getErrByCode(JSON.Get([]byte(res.String()), "code").ToUint())
	}
	return res.String(), nil
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
