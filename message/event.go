package message

// Mirai-api-http事件类型一览
// https://github.com/project-mirai/mirai-api-http/blob/master/EventType.md

/**
 * Bot Event
 */
const (
	// EventBotOnline Bot登录成功
	EventBotOnline = "BotOnlineEvent"
	// EventBotOfflineActive Bot主动离线
	EventBotOfflineActive = "BotOfflineEventActive"
	// EventBotOfflineForce Bot被挤下线
	EventBotOfflineForce = "BotOfflineEventForce"
	// EventBotOfflineDropped Bot被服务器断开或因网络问题而掉线
	EventBotOfflineDropped = "BotOfflineEventDropped"
	// EventBotRelogin Bot主动重新登录
	EventBotRelogin = "BotReloginEvent"
	// EventBotGroupPermissionChange Bot在群里的权限被改变. 操作人一定是群主
	EventBotGroupPermissionChange = "BotGroupPermissionChangeEvent"
	// EventBotMute Bot被禁言
	EventBotMute = "BotMuteEvent"
	// EventBotUnmute Bot被取消禁言
	EventBotUnmute = "BotUnmuteEvent"
	// EventBotJoinGroup Bot加入了一个新群
	EventBotJoinGroup = "BotJoinGroupEvent"
	// EventBotLeaveActive Bot主动退出一个群
	EventBotLeaveActive = "BotLeaveEventActive"
	// EventBotLeaveKick Bot被踢出一个群
	EventBotLeaveKick = "BotLeaveEventKick"
)

/**
 * Message Event Type
 */
const (
	// EventReceiveFriendMessage 好友消息
	EventReceiveFriendMessage = "FriendMessage"
	// EventReceiveGroupMessage 群组消息
	EventReceiveGroupMessage = "GroupMessage"
	// EventReceiveTempMessage 临时消息
	EventReceiveTempMessage = "TempMessage"
)

/**
 * Group Setting Event
 */

const (
	// EventGroupNameChange 某个群名改变
	EventGroupNameChange = "GroupNameChangeEvent"
	// EventGroupEntranceAnnouncementChange 某群入群公告改变
	EventGroupEntranceAnnouncementChange = "GroupEntranceAnnouncementChangeEvent"
	// EventGroupAllowAnonymousChat 匿名聊天
	EventGroupAllowAnonymousChat = "GroupAllowAnonymousChatEvent"
	// EventGroupAllowConfessTalk 坦白说
	EventGroupAllowConfessTalk = "GroupAllowConfessTalkEvent"
	// EventGroupAllowMemberInvite 允许群员邀请好友加群
	EventGroupAllowMemberInvite = "GroupAllowMemberInviteEvent"
)

/**
 * Group Member / Message Event
 */
const (
	// EventMemberJoinRequest 用户入群申请（Bot需要有管理员权限）
	EventMemberJoinRequest = "MemberJoinRequestEvent"
	// EventMemberJoin 新人入群的事件
	EventMemberJoin = "MemberJoinEvent"
	// EventMemberLeaveKick 成员被踢出群（该成员不是Bot）
	EventMemberLeaveKick = "MemberLeaveEventKick"
	// EventMemberLeaveQuit 成员主动离群（该成员不是Bot）
	EventMemberLeaveQuit = "MemberLeaveEventQuit"

	// EventGroupMuteAll 全员禁言
	EventGroupMuteAll = "GroupMuteAllEvent"
	// EventGroupRecall 群消息撤回
	EventGroupRecall = "GroupRecallEvent"
	// EventMemberMute 群成员被禁言事件（该成员不可能是Bot，见BotMuteEvent）
	EventMemberMute = "MemberMuteEvent"
	// EventMemberUnmute 群成员被取消禁言事件（该成员不可能是Bot，见BotUnmuteEvent）
	EventMemberUnmute = "MemberUnmuteEvent"

	// EventMemberCardChange 群名片改动
	EventMemberCardChange = "MemberCardChangeEvent"
	// EventMemberSpecialTitleChange 群头衔改动（只有群主有操作限权）
	EventMemberSpecialTitleChange = "MemberSpecialTitleChangeEvent"
	// EventMemberPermissionChange 成员权限改变的事件（该成员不可能是Bot，见BotGroupPermissionChangeEvent）
	EventMemberPermissionChange = "MemberPermissionChangeEvent"

	// EventBotInvitedJoinGroupRequest Bot被邀请入群申请
	EventBotInvitedJoinGroupRequest = "BotInvitedJoinGroupRequestEvent"
)

/**
 * Friend Event
 */
const (
	// EventNewFriendRequest 添加好友申请
	EventNewFriendRequest = "NewFriendRequestEvent"
	// EventFriendRecall 好友消息撤回
	EventFriendRecall = "FriendRecallEvent"
)

// Sender 消息发送者
type Sender struct {
	// ID QQ号
	ID uint `json:"id,omitempty"`
	// Group (GroupMessage)消息来源群信息
	Group Group `json:"group,omitempty"`
	// NickName (FriendMessage)发送者昵称
	NickName string `json:"nickname,omitempty"`
	// Remark (FriendMessage)发送者备注
	Remark string `json:"remark,omitempty"`

	// MemberName (GroupMessage)发送者群昵称
	MemberName string `json:"memberName,omitempty"`
	// Permission (GroupMessage)发送者在群中的角色
	Permission string `json:"permission,omitempty"`
}

// Member 成员(被操作对象)
type Member struct {
	Group      Group  `json:"group"`
	ID         int64  `json:"id"`
	MemberName string `json:"memberName"`
	Permission string `json:"permission"`
}

// Operator 操作者
type Operator struct {
	// ID QQ号
	ID uint `json:"id,omitempty"`
	// Group (GroupMessage)消息来源群信息
	Group Group `json:"group,omitempty"`
	// MemberName 群名片
	MemberName string `json:"memberName"`
	// Permission 权限
	Permission string `json:"permission"`
}

// Event 事件
type Event struct {
	// Type 事件类型
	Type string `json:"type"`
	// MessageChain (ReceiveMessage)消息链
	MessageChain []Message `json:"messageChain"`
	// Sender (ReceiveMessage)发送者信息
	Sender Sender `json:"sender"`
	// EventID 事件ID
	EventID uint `json:"eventId"`
	// FromID 操作人
	FromID uint `json:"fromId"`
	// GroupID 群号
	GroupID uint `json:"groupId"`
	// Operator 操作人
	Operator Operator `json:"operator"`
	// DurationSeconds 禁言时间
	DurationSeconds int64 `json:"durationSeconds"`
}

// Group QQ群
type Group struct {
	// ID 消息来源群号
	ID uint `json:"id,omitempty"`
	// Name 消息来源群名
	Name string `json:"name,omitempty"`
	// Permisson bot在群中的角色
	Permisson string `json:"permisson,omitempty"`
}
