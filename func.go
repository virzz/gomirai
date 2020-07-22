package gomirai

import "github.com/virzz/gomirai/message"

// SendGroupMessageWithBot 发送群消息
func SendGroupMessageWithBot(bot *Bot, qq, quote uint, msg ...message.Message) (uint, error) {
	return bot.SendGroupMessage(qq, quote, msg...)
}

// SendFriendMessageWithBot 发送好友消息
func SendFriendMessageWithBot(bot *Bot, group, quote uint, msg ...message.Message) (uint, error) {
	return bot.SendGroupMessage(group, quote, msg...)
}
