package gomirai

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"github.com/virzz/gomirai/message"
)

// Bot 对应一个机器人账号,进行所有对账号相关操作
type Bot struct {
	QQ          uint
	SessionKey  string
	Client      *Client
	Logger      *logrus.Entry
	fetchTime   time.Duration
	size        int
	currentSize int
	Chan        chan message.ComplexEvent // message.Event
}

// SetChannel Channel相关设置
func (b *Bot) SetChannel(time time.Duration, size int) {
	b.Chan = make(chan message.ComplexEvent, size)
	b.size = size
	b.currentSize = 0
	b.fetchTime = time
}

// SendFriendMessage 发送好友消息
// qq 好友qq
// quote 引用消息id 0为不引用
// msg 消息内容
func (b *Bot) SendFriendMessage(qq, quote uint, msg ...message.Message) (uint, error) {
	data := map[string]interface{}{"sessionKey": b.SessionKey, "qq": qq, "messageChain": msg}
	if quote != 0 {
		data["quote"] = quote
	}
	res, err := b.Client.doPost("/sendFriendMessage", data)
	if err != nil {
		return 0, err
	}
	b.Logger.Infoln("Send FriendMessage to", qq)
	return uint(gjson.Get(res, "messageId").Uint()), nil
}

// SendTempMessage 发送临时消息
// qq 好友qq
// group 群qq
// msg 消息内容
func (b *Bot) SendTempMessage(group, qq uint, msg ...message.Message) (uint, error) {
	data := map[string]interface{}{"sessionKey": b.SessionKey, "qq": qq, "group": group, "messageChain": msg}
	res, err := b.Client.doPost("/sendTempMessage", data)
	if err != nil {
		return 0, err
	}
	b.Logger.Infoln("Send TempMessage to ", qq)
	return uint(gjson.Get(res, "messageId").Uint()), nil
}

// SendGroupMessage 发送好友消息
// group 群qq
// quote 引用消息id 0为不引用
// msg 消息内容
func (b *Bot) SendGroupMessage(group, quote uint, msg ...message.Message) (uint, error) {
	data := map[string]interface{}{"sessionKey": b.SessionKey, "group": group, "messageChain": msg}
	if quote != 0 {
		data["quote"] = quote
	}
	res, err := b.Client.doPost("/sendGroupMessage", data)
	if err != nil {
		return 0, err
	}
	b.Logger.Infoln("Send FriendMessage to", group)
	return uint(gjson.Get(res, "messageId").Uint()), nil
}

// MuteGroupMember 禁言群成员
func (b *Bot) MuteGroupMember(group, qq uint, time int) error {
	data := map[string]interface{}{
		"sessionKey": b.SessionKey,
		"target":     group,
		"memberId":   qq,
		"time":       time,
	}
	_, err := b.Client.doPost("/mute", data)
	if err != nil {
		return err
	}
	b.Logger.Debugln("MuteGroupMember", qq, "in", group)
	return nil
}

// UnMuteGroupMember 取消禁言群成员
func (b *Bot) UnMuteGroupMember(group, qq uint) error {
	data := map[string]interface{}{
		"sessionKey": b.SessionKey,
		"target":     group,
		"memberId":   qq,
	}
	_, err := b.Client.doPost("/unmute", data)
	if err != nil {
		return err
	}
	b.Logger.Debugln("UnMuteGroupMember", qq, "in", group)
	return nil
}

// FetchMessages 获取消息
func (b *Bot) FetchMessages() error {
	t := time.NewTicker(b.fetchTime)
	for {
		res, err := b.Client.doGet("/fetchMessage", map[string]string{
			"sessionKey": b.SessionKey,
			"count":      strconv.Itoa(b.size),
		})
		if err != nil {
			return err
		}
		// b.Logger.Debugln(res)
		events := gjson.Get(res, "data").Array()
		for _, event := range events {
			if event.Get("type").String() == message.EventGroupMuteAll {
				continue
			}
			var c message.ComplexEvent
			if err := json.Unmarshal([]byte(event.Raw), &c); err != nil {
				b.Logger.Errorln("Unmarshal Event", err)
				continue
			}
			if len(b.Chan) == b.size {
				<-b.Chan
			}
			b.Chan <- c
		}
		// for _, event := range events {
		// 	if len(b.Chan) == b.size {
		// 		<-b.Chan
		// 	}
		// 	b.Chan <- event
		// }
		<-t.C
	}
}

const (
	// OperateAgree 同意入群
	OperateAgree = iota
	// OperateRefuse 拒绝入群
	OperateRefuse
	// OperateIgnore 忽略请求
	OperateIgnore
	// OperateRefuseBan 拒绝入群并添加黑名单，(腾讯)不再接收该用户的入群申请
	OperateRefuseBan
	// OperateIgnoreBan 忽略入群并添加黑名单，(腾讯)不再接收该用户的入群申请
	OperateIgnoreBan
)

// RespondMemberJoinRequest 响应用户加群请求
// operate	说明
// 0	同意入群
// 1	拒绝入群
// 2	忽略请求
// 3	拒绝入群并添加黑名单，不再接收该用户的入群申请
// 4	忽略入群并添加黑名单，不再接收该用户的入群申请
func (b *Bot) RespondMemberJoinRequest(eventID, fromID, groupID uint, operate int, message string) error {
	data := map[string]interface{}{"sessionKey": b.SessionKey, "eventId": eventID, "fromId": fromID, "groupId": groupID, "operate": operate, "message": message}
	_, err := b.Client.doPost("/resp/memberJoinRequestEvent", data)
	if err != nil {
		return err
	}
	b.Logger.Infoln("Respond Member Join Request ", fromID, " join ", groupID, " operate: ", operate)
	return nil
}

// SendImageMessage 使用此方法向指定对象（群或好友）发送图片消息
// 除非需要通过此手段获取imageId，否则不推荐使用该接口
// 请保证 qq group 不同时有值
func (b *Bot) SendImageMessage(qq, group int64, urls ...string) (imageIds []string, err error) {
	if qq*group == 0 {
		return nil, errors.New("非法参数")
	}
	data := map[string]interface{}{"sessionKey": b.SessionKey, "urls": urls}
	if qq == 0 {
		data["group"] = group
	} else {
		data["qq"] = qq
	}
	res, err := b.Client.doPost("/sendImageMessage", data)
	if err != nil {
		return nil, err
	}
	b.Logger.Info("Send Images")
	imageIds = make([]string, 0)
	for _, id := range gjson.Parse(res).Array() {
		imageIds = append(imageIds, id.String())
	}
	return
}

// UploadImage -
func (b *Bot) UploadImage(t string, imgFilepath string) (string, error) {
	imgReader, err := os.Open(imgFilepath)
	if err != nil {
		return "", err
	}
	defer imgReader.Close()

	data := map[string]interface{}{"sessionKey": b.SessionKey, "type": t, "img": imgReader}
	res, err := b.Client.doPostWithFormData("/uploadImage", data)
	if err != nil {
		return "", err
	}
	b.Logger.Info("UploadFriendImage ", imgFilepath)
	return gjson.Get(res, "imageId").String(), nil
}

// Recall 撤回消息
// target 消息id
func (b *Bot) Recall(target uint) error {
	data := map[string]interface{}{"sessionKey": b.SessionKey, "target": target}
	_, err := b.Client.doPost("/recall", data)
	return err
}

// RecallGroupMessage 撤回消息
func (b *Bot) RecallGroupMessage(sourceID uint) error {
	return b.Recall(sourceID)
}
