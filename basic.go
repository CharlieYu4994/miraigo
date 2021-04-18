package miraigo

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

func apiPostJSON(url string, data, res interface{}) error {
	tmp, _ := json.Marshal(data)

	resp, err := http.Post(url,
		"application/json", bytes.NewBuffer(tmp))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		return err
	}
	return nil
}

// 检查错误
func checkError(resp Response) error {
	switch resp.Code {
	case Success:
		return nil
	case WrongAuthKey:
		return errors.New("WrongAuthKey")
	case BotNotFound:
		return errors.New("BotNotFound")
	case SessionNotFound:
		return errors.New("SessionNotFound")
	case SessionNotActivated:
		return errors.New("SessionNotActived")
	case TargetNotFound:
		return errors.New("TargetNotFound")
	case FileNotFound:
		return errors.New("FileNotFound")
	case PermissionDenied:
		return errors.New("PermissionDenied")
	case Muted:
		return errors.New("Muted")
	case MessageTooLong:
		return errors.New("MessageTooLong")
	case WrongParams:
		return errors.New("WrongParams")
	default:
		return errors.New("UnknownError")
	}
}

func worker(event chan *Event, b *Bot) {
	for e := range event {
		for _, lookup := range b.lookupTable {
			if e.Type == lookup.typ {
				for _, m := range e.MessageChain {
					if lookup.matcher.MatchString(m.Text) {
						lookup.operate(b, e)
					}
				}
			}
		}
	}
}
