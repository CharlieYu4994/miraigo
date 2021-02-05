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
		return errors.New("Wrong AuthKey")
	case BotNotFound:
		return errors.New("Bot Not Found")
	case SessionNotFound:
		return errors.New("Session Not Found")
	case SessionNotActivated:
		return errors.New("Session Not Actived")
	case TargetNotFound:
		return errors.New("Target Not Found")
	case FileNotFound:
		return errors.New("File Not Found")
	case PermissionDenied:
		return errors.New("Permission Denied")
	case Muted:
		return errors.New("Muted")
	case MessageTooLong:
		return errors.New("Message Too Long")
	case WrongParams:
		return errors.New("Wrong Params")
	default:
		return errors.New("Unknown Error")
	}
}
