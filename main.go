package miraigo

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// NewBot 创建新的机器人实例
// @param  host    string "mirai 的地址(应包括协议名)"
//         authkey string "连接密钥"
//         qq      int64  "机器人的 QQ 号"
// @return *Bot  "机器人实例"
//         error "错误"
func NewBot(host, authkey string, qq int64) (*Bot, error) {
	session, err := getSession(host, authkey)
	if err != nil {
		return nil, err
	}
	err = activeSession(host, session, qq)
	if err != nil {
		return nil, err
	}
	err = setSession(host, session, 0, true)
	if err != nil {
		return nil, err
	}
	mws, mc, err := newListener(session, host, MsgEvent)
	if err != nil {
		return nil, err
	}
	ews, ec, err := newListener(session, host, EvnEvent)
	if err != nil {
		return nil, err
	}
	return &Bot{qq: qq, session: session, url: host,
		Message: mc, Event: ec, listeners: struct {
			msgListener   *WSListener
			eventListener *WSListener
		}{mws, ews}}, nil
}

// Close 关闭 Bot
// @return error
func (b *Bot) Close() error {
	tmp := Request{SessionKey: b.session, QQ: b.qq}
	var res Response

	err := apiPostJSON(b.url+Release, tmp, &res)
	if err != nil {
		return err
	}

	b.listeners.msgListener.stop()
	b.listeners.eventListener.stop()

	return checkError(res)
}

// AddEvent 注册事件
func (b *Bot) AddEvent(condition string, operate func(b Bot, e *Event)) {
	b.lookupTable = append(b.lookupTable, newLookup(condition, operate))
}

// Start 开始机器人任务
func (b *Bot) Start() error {
	b.listeners.msgListener.startListener()
	b.listeners.eventListener.startListener()
	var msg, evn *Event
	for {
		select {
		case msg = <-b.Message:
			go 
		case evn = <-b.Event:
			fmt.Printf("EVN: %#v\n", evn)
		}
	}
}

// 得到 SessionKey
func getSession(host, authkey string) (string, error) {
	tmp := Request{Authkey: authkey}
	var res Response

	err := apiPostJSON(host+Auth, tmp, &res)
	if err != nil {
		return "", err
	}

	err = checkError(res)
	if err != nil {
		return "", err
	}
	return res.Session, nil
}

// 激活 SessionKey
func activeSession(host, session string, qq int64) error {
	tmp := Request{SessionKey: session, QQ: qq}
	var res Response

	err := apiPostJSON(host+Verify, tmp, &res)
	if err != nil {
		return err
	}

	return checkError(res)
}

// 设置 Session
func setSession(host, session string, cache int32, websocket bool) error {
	tmp := Request{SessionKey: session, Websocket: true, CacheSize: cache}
	var res Response

	err := apiPostJSON(host+SessionConfig, tmp, &res)
	if err != nil {
		return err
	}

	return checkError(res)
}

func newLookup(c string, o func(b Bot, m *Event)) (*lookup, error) {
	tmp := strings.Split(c, ".")
	if len(tmp) <= 3 {
		return nil, errors.New("Wrong Format")
	}
	return &lookup{
		condition: map[string]string{
			"type": tmp[0],
			"id":   tmp[1],
			"msg":  tmp[2],
		},
		operate: o,
	}, nil
}
