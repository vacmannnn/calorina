package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tele "gopkg.in/telebot.v4"
)

type Dish struct {
	Name          string
	Calories      int
	Protein       int
	Fat           int
	Carbohydrates int
	Weight        int
}

type DailyEatingInfo struct {
	TotalCalories      int
	TotalProteins      int
	TotalFats          int
	TotalCarbohydrates int

	Date   string
	Dishes []Dish
}

func (dei DailyEatingInfo) String() string {
	pattern := `
Текущая дата: %s
Итого калорий за день: %d
Итого белков за день: %d
Итого жиров за день: %d
Итого углеводов за день: %d

Блюда съедены: %s`
	dishesNames := make([]string, len(dei.Dishes))
	for _, dish := range dei.Dishes {
		dishesNames = append(dishesNames, dish.Name)
	}

	return fmt.Sprintf(pattern, dei.Date, dei.TotalCalories, dei.TotalProteins,
		dei.TotalFats, dei.TotalCarbohydrates, strings.Join(dishesNames, " "))
}

var usersDishesInfo = map[int64]map[string]DailyEatingInfo{}

// todo:
// 1. sqlite to store users info
// 2. parse optional parameters (fat/protein/etc)
// 3. parse values as floats
// 4. move bot logic into separate function
// 5. deploy somewhere (optional)
// 6. readable readme
func main() {
	pref := tele.Settings{
		Token:  os.Getenv("TOKEN_CALORINA"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	// main message is:

	b.Handle(tele.OnText, func(c tele.Context) error {
		// All the text messages that weren't
		// captured by existing handlers.

		var (
			user = c.Sender()
			text = c.Text()
		)
		log.Println("got message from:", user.FirstName, "text:", text)

		if user.FirstName == "Arina" {
			return c.Send("<3")
		}

		log.Println(usersDishesInfo, usersDishesInfo[user.ID])
		curDayString := time.Now().Format("02.01.2006")
		if usersDishesInfo[user.ID] == nil {
			usersDishesInfo[user.ID] = map[string]DailyEatingInfo{}
		}
		curDay := usersDishesInfo[user.ID][curDayString]

		for _, dishString := range strings.Split(text, "\n") {
			dish := parseUsersDish(dishString)
			dailyEatingInfo := addDishToDay(curDay, dish)
			usersDishesInfo[user.ID][curDayString] = dailyEatingInfo
		}

		_, err := b.Send(user, usersDishesInfo[user.ID][curDayString].String())
		if err != nil {
			return err
		}

		// Instead, prefer a context short-hand:
		return nil
	})

	log.Println("starting telegram bot")
	b.Start()
}

func addDishToDay(dInfo DailyEatingInfo, d Dish) DailyEatingInfo {
	return DailyEatingInfo{
		Date:               time.Now().Format("02.01.2006"),
		TotalCalories:      dInfo.TotalCalories + d.Calories,
		TotalProteins:      dInfo.TotalProteins + d.Protein,
		TotalFats:          dInfo.TotalFats + d.Fat,
		TotalCarbohydrates: dInfo.TotalCarbohydrates + d.Carbohydrates,
		Dishes:             append(dInfo.Dishes, d),
	}
}

// название блюда / калориии / б / ж / у / грамм
// б-ж-у и граммы опциональны, можно не указывать их или пропустить через "-"
// если кол-во грамм не указано, то калории добавляются напрямую, иначе считаются по формуле "всего калорий" = "калории" * грамм * 0,01
func parseUsersDish(text string) Dish {
	// Expect input like: "DishName / 250 / 10 / 5 / 30 / 150"
	parts := strings.Split(text, " ")

	var d Dish
	if len(parts) == 0 {
		return d
	}

	// Name
	d.Name = parts[0]

	// helper to parse int, treating "-" or empty as zero and ignoring errors
	parseInt := func(s string) int {
		if s == "" || s == "-" {
			return 0
		}
		v, err := strconv.Atoi(s)
		if err != nil {
			return 0
		}
		return v
	}

	// calories
	if len(parts) >= 2 {
		d.Calories = parseInt(parts[1])
	}
	// protein
	if len(parts) >= 3 {
		d.Protein = parseInt(parts[2])
	}
	// fat
	if len(parts) >= 4 {
		d.Fat = parseInt(parts[3])
	}
	// carbs
	if len(parts) >= 5 {
		d.Carbohydrates = parseInt(parts[4])
	}
	// weight (grams)
	if len(parts) >= 6 {
		d.Weight = parseInt(parts[5])
	}

	// If weight specified (>0), scale nutrients and calories proportionally
	if d.Weight > 0 {
		factor := float64(d.Weight) * 0.01
		// scale
		dCalories := float64(d.Calories) * factor
		dProtein := float64(d.Protein) * factor
		dFat := float64(d.Fat) * factor
		dCarb := float64(d.Carbohydrates) * factor

		// assign back as ints
		d.Calories = int(dCalories + 0.5)
		d.Protein = int(dProtein + 0.5)
		d.Fat = int(dFat + 0.5)
		d.Carbohydrates = int(dCarb + 0.5)
	}

	return d
}
