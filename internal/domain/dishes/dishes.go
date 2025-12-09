package dishes

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Dish struct {
	Name          string
	Calories      int64
	Protein       int64
	Fat           int64
	Carbohydrates int64
	Weight        int64
}

// NewDish ожидает формат типа:
// название блюда / калориии / б / ж / у / грамм
// б-ж-у и граммы опциональны, можно не указывать их или пропустить через "-"
// если кол-во грамм не указано, то калории добавляются напрямую, иначе считаются по формуле "всего калорий" = "калории" * грамм * 0,01
func NewDish(text string) (Dish, bool) {
	var isTestingData bool
	parts := strings.Split(text, " ")

	var d Dish
	if len(parts) < 2 {
		return d, isTestingData
	}

	// Name
	d.Name = parts[0]
	d.Calories = parseFloatValueAsInt64(parts[1])

	switch len(parts) {
	case 7:
		isTestingData = parts[6] == "t"
		fallthrough
	case 6:
		d.Weight = parseFloatValueAsInt64(parts[5])
		fallthrough
	case 5:
		d.Carbohydrates = parseFloatValueAsInt64(parts[4])
		fallthrough
	case 4:
		d.Fat = parseFloatValueAsInt64(parts[3])
		fallthrough
	case 3:
		d.Protein = parseFloatValueAsInt64(parts[2])
	default:
		isTestingData = true
		return d, isTestingData
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

	return d, isTestingData
}

func parseFloatValueAsInt64(s string) int64 {
	if s == "-" {
		return 0
	}
	s = strings.ReplaceAll(s, ",", ".") // Replace comma with dot
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return int64(f)
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

func (dei DailyEatingInfo) AddDishToDay(d Dish) {
	dei.Date = time.Now().Format("02.01.2006")
	dei.TotalCalories += d.Calories
	dei.TotalProteins += d.Protein
	dei.TotalFats += d.Fat
	dei.TotalCarbohydrates += d.Carbohydrates
	dei.DishesNames = append(dei.DishesNames, d.Name)
}
