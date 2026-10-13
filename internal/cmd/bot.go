package main

import (
	"log"
	"strings"
	"time"

	repo "github.com/vacmannnn/calorina/internal/adapter/sqlite/dishes"
	"github.com/vacmannnn/calorina/internal/domain/dishes"
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

	_, err = b.botAPI.Send(user, b.dayFullString(user.ID, day))
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

	_, err = b.botAPI.Send(user, b.dayFullString(user.ID, day))
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

	_, err = b.botAPI.Send(user, b.dayString(user.ID, day))
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

	_, err = b.botAPI.Send(user, b.dayString(user.ID, day))
	if err != nil {
		log.Println(err)
	}

	return nil
}

func (b *Bot) helpMessage(c tele.Context) error {
	helpMessage := `/dishes - посмотреть информацию о блюдах за текущий день 
/delete ID - удалить блюдо за текущий день
/date - посмотреть результаты за конкретную дату (формат '19.01.2026')
/add - добавить блюдо к какому-то дню
/goal ККАЛ Б Ж У - задать цель на день
`

	return c.Send(helpMessage)
}

func (b *Bot) setGoal(c tele.Context) error {
	var (
		user = c.Sender()
		text = c.Text()
	)

	goal, err := b.service.SetGoal(user.ID, text)
	if err != nil {
		log.Println(err)
		return err
	}
	if goal == (dishes.Goal{}) {
		return c.Send("формат: /goal 2000 150 70 250")
	}

	return c.Send(goal.String())
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

	_, err = b.botAPI.Send(user, b.dayString(user.ID, day))
	if err != nil {
		log.Println(err)
	}

	return nil
}

func (b *Bot) dayString(userID int64, day dishes.DailyEatingInfo) string {
	return b.withGoal(userID, day, day.String())
}

func (b *Bot) dayFullString(userID int64, day dishes.DailyEatingInfo) string {
	return b.withGoal(userID, day, day.StringFullInfo())
}

func (b *Bot) withGoal(userID int64, day dishes.DailyEatingInfo, message string) string {
	goal, ok, err := b.service.GetGoal(userID)
	if err != nil || !ok {
		return message
	}

	return message + "\n" + day.RemainingString(goal)
}
