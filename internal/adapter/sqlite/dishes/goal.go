package dishes

import (
	"context"
	"database/sql"
	"errors"

	"github.com/vacmannnn/calorina/internal/adapter/sqlite/dishes/sqlc"
	"github.com/vacmannnn/calorina/internal/domain/dishes"
)

func (r *Repository) UpsertGoal(ctx context.Context, userID int64, goal dishes.Goal) error {
	return r.queries.UpsertGoal(ctx, sqlc.UpsertGoalParams{
		UserID:        userID,
		Calories:      goal.Calories,
		Proteins:      goal.Proteins,
		Fats:          goal.Fats,
		Carbohydrates: goal.Carbohydrates,
	})
}

func (r *Repository) GetGoal(ctx context.Context, userID int64) (dishes.Goal, bool, error) {
	goal, err := r.queries.GetGoal(ctx, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return dishes.Goal{}, false, nil
	}
	if err != nil {
		return dishes.Goal{}, false, err
	}

	return dishes.Goal{
		Calories:      goal.Calories,
		Proteins:      goal.Proteins,
		Fats:          goal.Fats,
		Carbohydrates: goal.Carbohydrates,
	}, true, nil
}
