package miraigo

import (
	"errors"
	"regexp"
	"strings"
)

// NewBot 创建新的机器人实例
// @param  host      string "mirai 的地址(应包括协议名)"
//         authkey   string "连接密钥"
//         qq        string "机器人的 QQ 号"
//         workerNum int    "工作线程数量"
// @return *Bot  "机器人实例"
//         error "错误"
//
func NewBot(host, authkey, qq string, workerNum int) (*Bot, error) {
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
	eventws, err := newListener(session, host, MsgEvent)
	if err != nil {
		return nil, err
	}

	return &Bot{qq: qq, session: session, url: host,
		eventListener: eventws, workerNum: workerNum}, nil
}

// Close 关闭 Bot
// @return error
//
func (b *Bot) Close() error {
	tmp := Request{SessionKey: b.session, QQ: b.qq}
	var res Response

	err := apiPostJSON(b.url+Release, tmp, &res)
	if err != nil {
		return err
	}

	b.eventListener.stop()

	close(b.eventListener.message)

	return checkError(res)
}

// AddEvent 注册事件
// @param condition string "type.msg"
//        authkey   string "操作函数"
//
func (b *Bot) AddEvent(condition string, operate Operate) {
	lookup, err := newLookup(condition, operate)
	if err != nil {
		panic(err)
	}
	b.lookupTable = append(b.lookupTable, lookup)
}

// Start 开始机器人任务
func (b *Bot) Start() error {
	err := b.eventListener.startListener()
	for i := 0; i <= b.workerNum; i++ {
		go worker(b.eventListener.message, b)
	}
	return err
}

// 得到 SessionKey
// @param  host    string "mirai 的地址(应包括协议名)"
//         authkey string "连接密钥"
// @return sessionKey string  "SessionKey"
//         error      error   "错误"
//
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
// @param  host       string "mirai 的地址(应包括协议名)"
//         sessionKey string  "SessionKey"
//         qq         string  "机器人的 QQ 号"
// @return error      error   "错误"
//
func activeSession(host, session, qq string) error {
	tmp := Request{SessionKey: session, QQ: qq}
	var res Response

	err := apiPostJSON(host+Verify, tmp, &res)
	if err != nil {
		return err
	}

	return checkError(res)
}

// 设置 Session
// @param  host       string "mirai 的地址(应包括协议名)"
//         sessionKey string  "SessionKey"
//         qq         string  "机器人的 QQ 号"
// @return error      error   "错误"
//
func setSession(host, session string, cache int, websocket bool) error {
	tmp := Request{SessionKey: session, Websocket: true, CacheSize: cache}
	var res Response

	err := apiPostJSON(host+SessionConfig, tmp, &res)
	if err != nil {
		return err
	}

	return checkError(res)
}

func newLookup(c string, o Operate) (*lookup, error) {
	tmp := strings.Split(c, ".")
	if len(tmp) <= 2 {
		return nil, errors.New("WrongFormat")
	}
	re := regexp.MustCompile(tmp[1])
	return &lookup{
		typ:     tmp[0],
		matcher: re,
		operate: o,
	}, nil
}
