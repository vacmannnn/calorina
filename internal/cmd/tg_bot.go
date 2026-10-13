package main

import (
	"os"
	"time"

	tele "gopkg.in/telebot.v4"
)

func initTgBotAPI() (*tele.Bot, error) {
	pref := tele.Settings{
		Token:  os.Getenv("TOKEN_CALORINA"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	return tele.NewBot(pref)
}
