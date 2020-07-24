package message

const (
	// MsgTypeSource -
	MsgTypeSource = "Source"
	// MsgTypeQuote 引用
	MsgTypeQuote = "Quote"
	// MsgTypeAt At
	MsgTypeAt = "At"
	// MsgTypeAtAll At所有人
	MsgTypeAtAll = "AtAll"
	// MsgTypeFace 表情
	MsgTypeFace = "Face"
	// MsgTypePlain 文本
	MsgTypePlain = "Plain"
	// MsgTypeImage 图片
	MsgTypeImage = "Image"
	// MsgTypeFlashImage 闪照
	MsgTypeFlashImage = "FlashImage"
	// MsgTypeXML XML
	MsgTypeXML = "Xml"
	// MsgTypeJSON JSON
	MsgTypeJSON = "Json"
	// MsgTypeApp App
	MsgTypeApp = "App"
	// MsgTypePoke 戳一戳
	MsgTypePoke = "Poke"
)

// Message 消息
type Message struct {
	Type string `json:"type,omitempty"`
	ID   uint   `json:"id,omitempty"`   //(Source,Quote)Source中表示消息id，Quote中表示被引用回复的原消息的id
	Time int64  `json:"time,omitempty"` //(Source) 发送时间

	GroupID  uint      `json:"groupId,omitempty"`  //(Quote)Quote中表示被引用回复的原消息的群号
	SenderID uint      `json:"senderId,omitempty"` //(Quote)Quote中表示被引用回复的原消息的发送者QQ号
	TargetID uint      `json:"targetId,omitempty"` //(Quote)Quote中表示被引用回复的原消息的接收者群号或QQ号
	Origin   []Message `json:"origin,omitempty"`   //(Quote)Quote中表示被引用回复的原消息的消息链对象

	Target  uint   `json:"target,omitempty"`  //(At)@的群员QQ号
	Display string `json:"display,omitempty"` //(At)@的显示文本

	FaceID int    `json:"faceId,omitempty"` //(Face)QQ表情的ID,发送时优先级比Name高
	Name   string `json:"name,omitempty"`   //(Face,Poke)Face中为QQ表情的拼音,Poke中为戳一戳的类型

	Text string `json:"text,omitempty"` //(Plain)纯文本

	ImageID   string `json:"imageId,omitempty"` //(Image,FlashImage)图片ID，注意消息类型，群图片和好友图片格式不一样，发送时优先级比ImageUrl高
	ImageURL  string `json:"url,omitempty"`     //(Image,FlashImage)图片url,发送时可使用网络图片的链接，优先级比ImagePath高；接收时为腾讯图片服务器的链接
	ImagePath string `json:"path,omitempty"`    //(Image,FlashImage)图片的路径，发送本地图片，相对路径于plugins/MiraiAPIHTTP/images

	XML     string `json:"xml,omitempty"`     //(Xml) xml消息本体
	JSON    string `json:"json,omitempty"`    //(Json) json消息本体
	Content string `json:"content,omitempty"` //(App) 不知道干嘛的，mirai也没有说明，估计是小程序连接？
}

// PlainMessage 文本消息
func PlainMessage(text string) Message {
	return Message{Type: MsgTypePlain, Text: text}
}

// AtMessage At消息
func AtMessage(target uint) Message {
	if target == 0 {
		return Message{Type: MsgTypeAtAll}
	}
	return Message{Type: MsgTypeAt, Target: target}
}

// FaceMessage 表情消息
func FaceMessage(faceID int) Message {
	return Message{Type: MsgTypeFace, FaceID: faceID}
}

// ImageMessage 图片消息
func ImageMessage(t, v string) Message {
	m := Message{Type: MsgTypeImage}
	switch t {
	case "id":
		m.ImageID = v
	case "url":
		m.ImageURL = v
	case "path":
		m.ImagePath = v
	default:
		return Message{}
	}
	return m
}

// FlashImageMessage 闪照消息
func FlashImageMessage(t, v string) Message {
	m := Message{Type: MsgTypeFlashImage}
	switch t {
	case "id":
		m.ImageID = v
	case "url":
		m.ImageURL = v
	case "path":
		m.ImagePath = v
	default:
		return Message{}
	}
	return m
}

// RichMessage 特殊消息
func RichMessage(t, content string) Message {
	m := Message{}
	switch t {
	case MsgTypeJSON:
		m.Type = t
		m.JSON = content
	case MsgTypeXML:
		m.Type = t
		m.XML = content
	case MsgTypeApp:
		m.Type = t
		m.Content = content
	default:
		return Message{}
	}
	return m
}

// PokeMessage 戳一戳消息
func PokeMessage(name string) Message {
	return Message{Type: MsgTypePoke, Name: name}
}
