package gomirai

import (
	"github.com/virzz/gomirai/message"
)

// SendGroupMessageWithBot 发送群消息
func SendGroupMessageWithBot(b *Bot, qq, quote int64, msg ...message.Message) (int64, error) {
	return b.SendGroupMessage(qq, quote, msg...)
}

// SendFriendMessageWithBot 发送好友消息
func SendFriendMessageWithBot(b *Bot, group, quote int64, msg ...message.Message) (int64, error) {
	return b.SendGroupMessage(group, quote, msg...)
}
