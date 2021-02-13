package miraigo

import "testing"

func TestMain(t *testing.T) {
	bot, _ := NewBot("http://192.168.1.2:8080",
		"s6RgGZuLCpMJKy17aEKW", 1472833065)
	bot.Start()
	bot.Close()
}
