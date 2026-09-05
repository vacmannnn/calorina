package dishes

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Dish struct {
	ID            int64
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
		isTestingData = true
		return d, isTestingData
	}
	d.Name, parts = parseDishName(parts)

	if len(parts) < 2 {
		isTestingData = true
		return d, isTestingData
	}
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

func parseDishName(parts []string) (string, []string) {
	var wordsCount int
	for i := range parts {
		if parseFloatValueAsInt64(parts[i]) != 0 {
			wordsCount = i
			break
		}
	}

	if wordsCount == 0 {
		return "", []string{}
	}
	return strings.Join(parts[:wordsCount], " "), parts[wordsCount-1:]
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

	Date   string
	Dishes []Dish
}

type Goal struct {
	Calories      int64
	Proteins      int64
	Fats          int64
	Carbohydrates int64
}

func NewGoal(text string) (Goal, bool) {
	parts := strings.Fields(text)
	if len(parts) != 5 {
		return Goal{}, false
	}

	goal := Goal{
		Calories:      parseFloatValueAsInt64(parts[1]),
		Proteins:      parseFloatValueAsInt64(parts[2]),
		Fats:          parseFloatValueAsInt64(parts[3]),
		Carbohydrates: parseFloatValueAsInt64(parts[4]),
	}
	if goal.Calories <= 0 || goal.Proteins < 0 || goal.Fats < 0 || goal.Carbohydrates < 0 {
		return Goal{}, false
	}

	return goal, true
}

func (g Goal) String() string {
	return fmt.Sprintf("цель на день: %d ккал, %d белков, %d жиров, %d углеводов",
		g.Calories, g.Proteins, g.Fats, g.Carbohydrates)
}

func (dei *DailyEatingInfo) RemainingString(goal Goal) string {
	return fmt.Sprintf(
		"остаток на сегодня: %d ккал, %d белков, %d жиров, %d углеводов",
		goal.Calories-dei.TotalCalories,
		goal.Proteins-dei.TotalProteins,
		goal.Fats-dei.TotalFats,
		goal.Carbohydrates-dei.TotalCarbohydrates,
	)
}

func (dei *DailyEatingInfo) String() string {
	pattern := `
Блюда съедены: %s`

	var dishesName string
	for _, dish := range dei.Dishes {
		if dish.Name == "" {
			continue
		}
		dishesName += "'" + strings.ToLower(dish.Name) + "', "
	}

	return dei.StringBasicInfo() + fmt.Sprintf(pattern, dishesName[:max(0, len(dishesName)-2)])
}

func (dei *DailyEatingInfo) StringFullInfo() string {
	//- Белков - %d гр.
	//- Жиров - %d гр.
	//- Углеводов - %d гр.
	//- Вес - %d гр.
	patternDishes := `
Блюдо %d:
- Название - %s
- Калорийность - %d ккал.
- id блюда - %d
`
	var output string
	for i, dish := range dei.Dishes {
		output += fmt.Sprintf(patternDishes, i+1, dish.Name, dish.Calories, dish.ID)
	}
	output = dei.StringBasicInfo() + output

	return output
}

func (dei *DailyEatingInfo) StringBasicInfo() string {
	patternDay := `
Текущая дата: %s
Итого калорий за день: %d
Итого белков за день: %d
Итого жиров за день: %d
Итого углеводов за день: %d
`

	return fmt.Sprintf(patternDay, dei.Date, dei.TotalCalories, dei.TotalProteins,
		dei.TotalFats, dei.TotalCarbohydrates)
}

func (dei *DailyEatingInfo) AddDishToDay(d Dish) {
	if dei.Date == "" {
		dei.Date = time.Now().Format("02.01.2006")
	}
	dei.TotalCalories += d.Calories
	dei.TotalProteins += d.Protein
	dei.TotalFats += d.Fat
	dei.TotalCarbohydrates += d.Carbohydrates
	dei.Dishes = append(dei.Dishes, d)
}

func (dei *DailyEatingInfo) RemoveDishByID(dishID int64) {
	for i, dish := range dei.Dishes {
		if dish.ID == dishID {
			dei.Dishes = append(dei.Dishes[:i], dei.Dishes[i+1:]...)
			dei.TotalCalories = max(dei.TotalCalories-dish.Calories, 0)
			dei.TotalProteins = max(dei.TotalProteins-dish.Protein, 0)
			dei.TotalFats = max(dei.TotalFats-dish.Fat, 0)
			dei.TotalCarbohydrates = max(dei.TotalCarbohydrates-dish.Carbohydrates, 0)
			return
		}
	}
}
