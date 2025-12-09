package main

import (
	"context"
	"log"
	"strings"
	"time"

	repo "github.com/vacmannnn/calorina/internal/adapter/sqlite/dishes"
	"github.com/vacmannnn/calorina/internal/domain/dishes"
	tele "gopkg.in/telebot.v4"
)

type Bot struct {
	repo   *repo.Repository
	botAPI *tele.Bot
}

func NewBot(repo *repo.Repository, botAPI *tele.Bot) *Bot {
	return &Bot{repo: repo, botAPI: botAPI}
}

func (b *Bot) handleText(c tele.Context) error {
	var (
		user = c.Sender()
		text = c.Text()
	)
	log.Println("got message from:", user.FirstName, "text:", text)

	curDayString := time.Now().Format("02.01.2006")
	day, err := b.repo.GetDailyInfo(context.TODO(), user.ID, curDayString)
	if err != nil {
		log.Println(err)
		return err
	}

	for _, dishString := range strings.Split(text, "\n") {
		dish, isTestData := dishes.NewDish(dishString)
		if !isTestData {
			err = b.repo.InsertDish(context.TODO(), dish)
			if err != nil {
				log.Println(err)
			}
		}

		day.AddDishToDay(dish)
		if !isTestData {
			err = b.repo.InsertDailyInfo(context.TODO(), user.ID, day)
			if err != nil {
				log.Println(err)
			}
		}
	}

	_, err = b.botAPI.Send(user, day.String())
	if err != nil {
		log.Println(err)
	}

	return nil
}
