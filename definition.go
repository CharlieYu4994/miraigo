package miraigo

import "regexp"

// 状态码定义
const (
	Success = iota
	WrongAuthKey
	BotNotFound
	SessionNotFound
	SessionNotActivated
	TargetNotFound
	FileNotFound
	PermissionDenied = 10
	Muted            = 20
	MessageTooLong   = 30
	WrongParams      = 400
)

// 定义请求地址
const (
	Auth          = "/auth"
	Verify        = "/verify"
	SessionConfig = "/config"
	Release       = "/release"
	MsgEvent      = "/message"
	EvnEvent      = "/event"
	AllEvent      = "/all"
)

const ()

// Bot 机器人对象
// 定义了机器人的基本参数
type Bot struct {
	workerNum     int
	qq            string
	session       string
	url           string
	lookupTable   []*lookup
	eventListener *WSListener
}

type lookup struct {
	typ     string // type GroupMessage id 2291598823 msg test
	matcher *regexp.Regexp
	operate func(b *Bot, m *Event)
}

// WSListener Websocket 监听器
type WSListener struct {
	url     string
	origin  string
	message chan *Event
	quit    chan bool
}

// Response 消息类返回信息模板
// 定义了返回的信息与结构
type Response struct {
	Code         int         `json:"code,omitempty"`
	MessageID    int         `json:"messageId,omitempty"`
	Msg          string      `json:"msg,omitempty"`
	ErrorMessage string      `json:"errorMessage,omitempty"`
	Session      string      `json:"session,omitempty"`
	ImageID      string      `json:"imageId,omitempty"`
	VoiceID      string      `json:"voiceId,omitempty"`
	URL          string      `json:"url,omitempty"`
	Path         string      `json:"path,omitempty"`
	Data         interface{} `json:"data,omitempty"`
}

// Request 请求模板
// 定义了一次请求需要的数据与结构
type Request struct {
	Websocket    bool     `json:"enableWebsocket"`
	Quote        int      `json:"quote,omitempty"`
	CacheSize    int      `json:"cacheSize,omitempty"`
	Target       string   `json:"target,omitempty"`
	Group        string   `json:"group,omitempty"`
	QQ           string   `json:"qq,omitempty"`
	Authkey      string   `json:"authKey,omitempty"`
	SessionKey   string   `json:"sessionKey,omitempty"`
	URLs         []string `json:"urls,omitempty"`
	MessageChain `json:"messageChain,omitempty"`
}

// Message 消息模板
// 定义了消息的组成
type Message struct {
	FaceID       int    `json:"faceId,omitempty"`
	Time         int    `json:"time,omitempty"`
	GroupID      string `json:"groupId,omitempty"`
	ID           string `json:"id,omitempty"`
	SenderID     string `json:"senderId,omitempty"`
	Target       string `json:"target,omitempty"`
	TargetID     string `json:"targetId,omitempty"`
	Content      string `json:"content,omitempty"`
	Display      string `json:"display,omitempty"`
	ImageID      string `json:"imageId,omitempty"`
	JSON         string `json:"json,omitempty"`
	Name         string `json:"name,omitempty"`
	Path         string `json:"path,omitempty"`
	Text         string `json:"text,omitempty"`
	Type         string `json:"type,omitempty"`
	URL          string `json:"url,omitempty"`
	VoiceID      string `json:"voiceId,omitempty"`
	XML          string `json:"xml,omitempty"`
	MessageChain `json:"origin,omitempty"`
}

// Persion QQ 中的人
// 定义了 QQ 中的人的结构
type Persion struct {
	ID           string `json:"id,omitempty"`
	Nickname     string `json:"nickname,omitempty"`
	Remaek       string `json:"remark,omitempty"`
	MemberName   string `json:"memberName,omitempty"`
	SpecialTitle string `json:"specialTitle,omitempty"`
	Group        `json:"group,omitempty"`
}

// Group 群
// 定义了群的结构
type Group struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Permission string `json:"permission,omitempty"`
}

// GroupConfig 群设置
// 定义了群设置
type GroupConfig struct {
	ConfessTalk       bool   `json:"confessTalk"`
	AllowMemberInvite bool   `json:"allowMemberInvite"`
	AutoApprove       bool   `json:"autoApprove"`
	AnonymousChat     bool   `json:"anonymousChat"`
	Name              string `json:"name,omitempty"`
	Announcement      string `json:"announcement,omitempty"`
}

// Event 事件模型
// 定义了在 websocket 接收的哦的 event
type Event struct {
	Type         string  `json:"type,omitempty"`
	Sender       Persion `json:"sender,omitempty"`
	MessageChain `json:"messageChain,omitempty"`
}

// PersionList 人员列表
// 定义了 QQ 中人员列表的结构
// 此列表为好友列表与群成员列表复用
type PersionList []*Persion

// GroupList 群列表
// 定义了好友列表的结构
type GroupList []*Group

// MessageChain 消息链
// 定义了消息链的结构
type MessageChain []*Message

// ImageChain 图片 ID 列表
// 定义了图片 ID 的响应结构
type ImageChain []string

type Operate func(b *Bot, m *Event)
