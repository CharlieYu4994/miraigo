package miraigo

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
	return &Bot{qq: qq, session: session, url: host}, nil
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

	return checkError(res)
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
