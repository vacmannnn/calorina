package main

import (
	"log"

	"github.com/vacmannnn/calorina/internal/service/dishes"
	tele "gopkg.in/telebot.v4"
)

func logMessage(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		log.Println("got message from:", c.Sender().FirstName, "text:", c.Text())
		return next(c)
	}
}

// todo:
// 5. deploy somewhere (optional)
// 6. readable readme
// 12. linter fixes
// 13. add buttons to remove and help
// 15. error handling
func main() {
	tgBotAPI, err := initTgBotAPI()
	if err != nil {
		log.Fatal(err)
	}

	repos, closer := initDB()
	defer closer()

	dishesSerivce := dishes.NewService(repos)

	bot := NewBot(dishesSerivce, tgBotAPI)

	tgBotAPI.Handle("/dishes", bot.printDishes, logMessage)
	tgBotAPI.Handle("/delete", bot.deleteDishes, logMessage)
	tgBotAPI.Handle("/date", bot.checkDate, logMessage)
	tgBotAPI.Handle("/add", bot.addToDate, logMessage)
	tgBotAPI.Handle("/help", bot.helpMessage, logMessage)
	tgBotAPI.Handle(tele.OnText, bot.handleText, logMessage)

	log.Println("starting telegram bot")
	tgBotAPI.Start()
}
