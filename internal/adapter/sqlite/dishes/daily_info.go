package dishes

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/vacmannnn/calorina/internal/adapter/sqlite/dishes/sqlc"
	"github.com/vacmannnn/calorina/internal/domain/dishes"
)

func (r *Repository) GetDailyInfo(ctx context.Context, userID int64, date string) (dishes.DailyEatingInfo, error) {
	info, err := r.queries.GetByUserIDAndDate(ctx, sqlc.GetByUserIDAndDateParams{
		UserID: userID,
		Date:   date,
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return dishes.DailyEatingInfo{}, err
	}

	var idInts []int64
	for _, ids := range strings.Split(strings.TrimSpace(info.DishesIds), ";") {
		id, err := strconv.Atoi(ids)
		if err != nil {
			log.Println("failed to parse dish id", ids)
		}
		idInts = append(idInts, int64(id))
	}
	ds, err := r.GetDishes(ctx, idInts)
	if err != nil {
		log.Println("failed to get dishes", idInts)
	}

	return dishes.DailyEatingInfo{
		TotalCalories:      info.TotalCalories,
		TotalProteins:      info.TotalProteins.Int64,
		TotalFats:          info.TotalFats.Int64,
		TotalCarbohydrates: info.TotalCarbohydrates.Int64,
		Date:               info.Date,
		Dishes:             ds,
	}, nil
}

func (r *Repository) InsertDailyInfo(ctx context.Context, userID int64, d dishes.DailyEatingInfo) error {
	dishesIDs := make([]string, 0, len(d.Dishes))
	for _, dish := range d.Dishes {
		dishesIDs = append(dishesIDs, strconv.Itoa(int(dish.ID)))
	}

	return r.queries.UpsertDailyInfo(ctx, sqlc.UpsertDailyInfoParams{
		UserID:        userID,
		Date:          d.Date,
		TotalCalories: d.TotalCalories,
		TotalProteins: sql.NullInt64{
			Int64: d.TotalProteins,
			Valid: true,
		},
		TotalFats: sql.NullInt64{
			Int64: d.TotalFats,
			Valid: true,
		},
		TotalCarbohydrates: sql.NullInt64{
			Int64: d.TotalCarbohydrates,
			Valid: true,
		},
		DishesIds: strings.Join(dishesIDs, ";"),
	})
}
