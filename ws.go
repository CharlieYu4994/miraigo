package miraigo

import (
	"encoding/json"
	"fmt"
	"net/url"

	"golang.org/x/net/websocket"
)

// NewListener 创建监听器对象，并把管道赋值给 Bot
func (b *Bot) NewListener() (*WSListener, error) {
	tmp, err := url.Parse(b.url)
	wsURL := "ws://" + tmp.Host + Websocket + "?sessionKey=" + b.session
	if err != nil {
		return nil, err
	}

	c := make(chan *Message, 128)
	b.Message = c
	q := make(chan bool, 1)

	return &WSListener{wsURL, b.url, c, q}, nil
}

// StartListener 启动监听器
func (w *WSListener) StartListener() error {
	var msg = make([]byte, 16*1024)
	ws, err := websocket.Dial(w.url, "", w.origin)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-w.quit:
				return
			default:
				n, err := ws.Read(msg)
				if err != nil {
					fmt.Println(err)
				}
				var tmp Message
				err = json.Unmarshal(msg[:n], &tmp)
				if err != nil {
					fmt.Println(err)
				}
				w.message <- &tmp
			}
		}
	}()
	return nil
}
