package dishes

import (
	"context"
	"database/sql"

	"github.com/vacmannnn/calorina/internal/adapter/sqlite/dishes/sqlc"
	"github.com/vacmannnn/calorina/internal/domain/dishes"
)

func (r *Repository) InsertDish(ctx context.Context, d dishes.Dish) error {
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

func (r *Repository) GetDish(ctx context.Context, dishIDs []int64) ([]dishes.Dish, error) {
	results := make([]dishes.Dish, 0, len(dishIDs))

	for _, name := range dishIDs {
		dishName, err := r.queries.GetDishNameByID(ctx, name)
		if err != nil {
			continue
		}

		results = append(results, dishes.Dish{
			Name: dishName,
		})
	}

	return results, nil
}
