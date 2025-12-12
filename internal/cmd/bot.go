package main

import (
	"log"
	"strings"
	"time"

	repo "github.com/vacmannnn/calorina/internal/adapter/sqlite/dishes"
	dservice "github.com/vacmannnn/calorina/internal/service/dishes"
	tele "gopkg.in/telebot.v4"
)

type Bot struct {
	repo    *repo.Repository
	botAPI  *tele.Bot
	service *dservice.Service
}

func NewBot(service *dservice.Service, botAPI *tele.Bot) *Bot {
	return &Bot{
		service: service,
		botAPI:  botAPI,
	}
}

func (b *Bot) printDishes(c tele.Context) error {
	var (
		user = c.Sender()
	)

	day, err := b.service.GetSpecificDayInfo(user.ID, time.Now().Format("02.01.2006"))
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = b.botAPI.Send(user, day.StringFullInfo())
	if err != nil {
		log.Println(err)
	}

	return nil
}

func (b *Bot) checkDate(c tele.Context) error {
	var (
		user = c.Sender()
		text = c.Text()
	)

	sp := strings.Split(text, " ")
	if len(sp) != 2 {
		return c.Send("better luck next time")
	}

	curDayString := sp[1]
	day, err := b.service.GetSpecificDayInfo(user.ID, curDayString)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = b.botAPI.Send(user, day.StringFullInfo())
	if err != nil {
		log.Println(err)
	}

	return nil
}

func (b *Bot) deleteDishes(c tele.Context) error {
	var (
		user = c.Sender()
		text = c.Text()
	)

	curDayString := time.Now().Format("02.01.2006")
	day, err := b.service.RemoveDishesFromDay(user.ID, text, curDayString)
	if err != nil {
		log.Println(err)
	}

	_, err = b.botAPI.Send(user, day.String())
	if err != nil {
		log.Println(err)
	}
	return err
}

func (b *Bot) addToDate(c tele.Context) error {
	var (
		user = c.Sender()
		text = c.Text()
	)

	lines := strings.Split(text, "\n")
	if len(strings.Split(lines[0], " ")) != 2 {
		return c.Send("better luck next time")
	}

	curDayString := strings.Split(lines[0], " ")[1]
	day, err := b.service.AddDishesToDay(user.ID, strings.Join(lines[1:], "\n"), curDayString)
	if err != nil {
		log.Println(err)
	}

	_, err = b.botAPI.Send(user, day.String())
	if err != nil {
		log.Println(err)
	}

	return nil
}

func (b *Bot) handleText(c tele.Context) error {
	var (
		user = c.Sender()
		text = c.Text()
	)

	curDayString := time.Now().Format("02.01.2006")
	day, err := b.service.AddDishesToDay(user.ID, text, curDayString)
	if err != nil {
		log.Println(err)
	}

	_, err = b.botAPI.Send(user, day.String())
	if err != nil {
		log.Println(err)
	}

	return nil
}
