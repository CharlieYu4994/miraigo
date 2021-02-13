package miraigo

import (
	"encoding/json"
	"fmt"
	"net/url"

	"golang.org/x/net/websocket"
)

// newListener 创建监听器对象
func newListener(session, origin, path string) (*WSListener, chan *Event, error) {
	tmp, err := url.Parse(origin)
	wsURL := "ws://" + tmp.Host + path + "?sessionKey=" + session
	if err != nil {
		return nil, nil, err
	}

	c := make(chan *Event, 128)
	q := make(chan bool, 1)

	return &WSListener{wsURL, origin, c, q}, c, nil
}

// startListener 启动监听器
func (w *WSListener) startListener() error {
	var msg = make([]byte, 16*1024)
	ws, err := websocket.Dial(w.url, "", w.origin)
	if err != nil {
		return err
	}

	go func() {
	MAINLOOP:
		for {
			select {
			case <-w.quit:
				break MAINLOOP
			default:
				n, err := ws.Read(msg)
				if err != nil {
					fmt.Println(err)
				}
				var tmp Event
				err = json.Unmarshal(msg[:n], &tmp)
				if err != nil {
					fmt.Println(err)
				}
				w.message <- &tmp
			}
		}
		ws.Close()
	}()
	return nil
}

func (w *WSListener) stop() {
	w.quit <- true
}
