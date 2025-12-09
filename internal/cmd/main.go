package main

import (
	"log"

	tele "gopkg.in/telebot.v4"
)

// todo:
// 10. remove dishes
// 11. check specific day
// 5. deploy somewhere (optional)
// 6. readable readme
// 12. linter fixes
func main() {
	tgBotAPI, err := initTgBotAPI()
	if err != nil {
		log.Fatal(err)
	}
	repos, closer := initDB()
	defer closer()
	bot := NewBot(repos, tgBotAPI)

	tgBotAPI.Handle(tele.OnText, bot.handleText)

	log.Println("starting telegram bot")
	tgBotAPI.Start()
}
