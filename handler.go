package gomirai

import (
	"github.com/virzz/gomirai/message"
)

// EventHandler -
type EventHandler struct {
	privateMessageHandlers []func(bot *Bot, chain message.Chain, sender message.Sender)
	groupMessageHandlers   []func(bot *Bot, chain message.Chain, sender message.Sender)
}
