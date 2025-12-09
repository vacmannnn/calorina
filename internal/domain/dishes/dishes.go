package dishes

import (
	"fmt"
	"strings"
)

type Dish struct {
	Name          string
	Calories      int64
	Protein       int64
	Fat           int64
	Carbohydrates int64
	Weight        int64
}

type DailyEatingInfo struct {
	TotalCalories      int64
	TotalProteins      int64
	TotalFats          int64
	TotalCarbohydrates int64

	Date        string
	DishesNames []string
}

func (dei DailyEatingInfo) String() string {
	pattern := `
Текущая дата: %s
Итого калорий за день: %d
Итого белков за день: %d
Итого жиров за день: %d
Итого углеводов за день: %d

Блюда съедены: %s`

	return fmt.Sprintf(pattern, dei.Date, dei.TotalCalories, dei.TotalProteins,
		dei.TotalFats, dei.TotalCarbohydrates, strings.Join(dei.DishesNames, " "))
}
