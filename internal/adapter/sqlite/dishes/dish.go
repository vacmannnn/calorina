package dishes

import (
	"context"
	"database/sql"

	"github.com/vacmannnn/calorina/internal/adapter/sqlite/dishes/sqlc"
	"github.com/vacmannnn/calorina/internal/domain/dishes"
)

func (r *Repository) InsertDish(ctx context.Context, d dishes.Dish) (int64, error) {
	return r.queries.UpsertDishInfo(ctx, sqlc.UpsertDishInfoParams{
		Name:     d.Name,
		Calories: d.Calories,
		Protein: sql.NullInt64{
			Valid: true,
			Int64: d.Protein,
		},
		Fat: sql.NullInt64{
			Valid: true,
			Int64: d.Fat,
		},
		Carbohydrates: sql.NullInt64{
			Valid: true,
			Int64: d.Carbohydrates,
		},
		Weight: sql.NullInt64{
			Valid: true,
			Int64: d.Weight,
		},
	})
}

func (r *Repository) GetDishes(ctx context.Context, dishIDs []int64) ([]dishes.Dish, error) {
	results := make([]dishes.Dish, 0, len(dishIDs))

	for _, name := range dishIDs {
		dish, err := r.queries.GetDishNameByID(ctx, name)
		if err != nil {
			continue
		}

		results = append(results, dishes.Dish{
			ID:            dish.ID,
			Name:          dish.Name,
			Calories:      dish.Calories,
			Protein:       dish.Protein.Int64,
			Fat:           dish.Fat.Int64,
			Carbohydrates: dish.Carbohydrates.Int64,
			Weight:        dish.Weight.Int64,
		})
	}

	return results, nil
}
