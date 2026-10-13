package dishes

import "testing"

func TestNewGoal(t *testing.T) {
	goal, ok := NewGoal("/goal 2000 150 70 250")
	if !ok {
		t.Fatal("expected goal to be parsed")
	}

	if goal.Calories != 2000 || goal.Proteins != 150 || goal.Fats != 70 || goal.Carbohydrates != 250 {
		t.Fatalf("unexpected goal: %#v", goal)
	}
}

func TestDailyEatingInfoRemainingString(t *testing.T) {
	day := DailyEatingInfo{
		TotalCalories:      500,
		TotalProteins:      40,
		TotalFats:          20,
		TotalCarbohydrates: 60,
	}
	goal := Goal{
		Calories:      2000,
		Proteins:      150,
		Fats:          70,
		Carbohydrates: 250,
	}

	got := day.RemainingString(goal)
	want := "остаток на сегодня: 1500 ккал, 110 белков, 50 жиров, 190 углеводов"
	if got != want {
		t.Fatalf("unexpected remaining string: got %q, want %q", got, want)
	}
}
