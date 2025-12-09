package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	repo "github.com/vacmannnn/calorina/internal/adapter/sqlite/dishes"
	"github.com/vacmannnn/calorina/internal/domain/dishes"
	tele "gopkg.in/telebot.v4"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// todo:
// 3. parse values as floats
// 4. move bot logic into separate function
// 5. deploy somewhere (optional)
// 6. readable readme
// 7. global refactoring
// 8. option to get daily info
// 9. multi-words dishes names
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

	db, err := sql.Open("sqlite3", "file:mydb.db?cache=shared&mode=rwc")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Apply migrations
	m, err := migrate.New(
		"file://schema/sqlite", // Path to your migration files
		"sqlite3://mydb.db") // Database URL
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully!")

	repos := repo.NewRepository(db)

	b.Handle(tele.OnText, func(c tele.Context) error {
		// All the text messages that weren't
		// captured by existing handlers.

		var (
			user = c.Sender()
			text = c.Text()
		)
		log.Println("got message from:", user.FirstName, "text:", text)

		curDayString := time.Now().Format("02.01.2006")
		day, err := repos.GetDailyInfo(context.TODO(), user.ID, curDayString)

		if len(strings.Split(text, " ")) == 1 {
			_, err = b.Send(user, day.String())
			return err
		}

		var dailyEatingInfo dishes.DailyEatingInfo
		for _, dishString := range strings.Split(text, "\n") {
			dish := parseUsersDish(dishString)
			err = repos.InsertDish(context.TODO(), dish)
			if err != nil {
				return err
			}
			dailyEatingInfo = addDishToDay(day, dish)
			err = repos.InsertDailyInfo(context.TODO(), user.ID, dailyEatingInfo)
			if err != nil {
				log.Println(err)
			}
		}

		_, err = b.Send(user, dailyEatingInfo.String())
		if err != nil {
			return err
		}

		// Instead, prefer a context short-hand:
		return nil
	})

	log.Println("starting telegram bot")
	b.Start()
}

func addDishToDay(dInfo dishes.DailyEatingInfo, d dishes.Dish) dishes.DailyEatingInfo {
	return dishes.DailyEatingInfo{
		Date:               time.Now().Format("02.01.2006"),
		TotalCalories:      dInfo.TotalCalories + d.Calories,
		TotalProteins:      dInfo.TotalProteins + d.Protein,
		TotalFats:          dInfo.TotalFats + d.Fat,
		TotalCarbohydrates: dInfo.TotalCarbohydrates + d.Carbohydrates,
		DishesNames:        append(dInfo.DishesNames, d.Name),
	}
}

// название блюда / калориии / б / ж / у / грамм
// б-ж-у и граммы опциональны, можно не указывать их или пропустить через "-"
// если кол-во грамм не указано, то калории добавляются напрямую, иначе считаются по формуле "всего калорий" = "калории" * грамм * 0,01
func parseUsersDish(text string) dishes.Dish {
	// Expect input like: "DishName / 250 / 10 / 5 / 30 / 150"
	parts := strings.Split(text, " ")

	var d dishes.Dish
	if len(parts) == 0 {
		return d
	}

	// Name
	d.Name = parts[0]

	// helper to parse int, treating "-" or empty as zero and ignoring errors
	parseInt := func(s string) int64 {
		if s == "" || s == "-" {
			return 0
		}
		v, err := strconv.Atoi(s)
		if err != nil {
			return 0
		}
		return int64(v)
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
		d.Calories = int64(dCalories + 0.5)
		d.Protein = int64(dProtein + 0.5)
		d.Fat = int64(dFat + 0.5)
		d.Carbohydrates = int64(dCarb + 0.5)
	}

	return d
}
