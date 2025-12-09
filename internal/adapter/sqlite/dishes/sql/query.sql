-- name: GetByUserIDAndDate :one
SELECT user_id, date, total_calories, total_proteins, total_fats, total_carbohydrates, dishes_ids
FROM daily_eating_info
WHERE user_id = @user_id AND date = @date;

-- name: UpsertDailyInfo :exec
INSERT INTO daily_eating_info (user_id, date, total_calories, total_proteins, total_fats, total_carbohydrates, dishes_ids)
VALUES (@user_id, @date, @total_calories, @total_proteins, @total_fats, @total_carbohydrates, @dishes_ids)
    ON CONFLICT(user_id, date) DO UPDATE SET
    total_calories = EXCLUDED.total_calories,
                                      total_proteins = EXCLUDED.total_proteins,
                                      total_fats = EXCLUDED.total_fats,
                                      total_carbohydrates = EXCLUDED.total_carbohydrates,
                                      dishes_ids = EXCLUDED.dishes_ids;

-- name: UpsertDishInfo :one
INSERT INTO dishes (name, calories, protein, fat, carbohydrates, weight)
VALUES (@name, @calories, @protein, @fat, @carbohydrates, @weight)
RETURNING id;

-- name: GetDishNameByID :one
SELECT *
FROM dishes d
WHERE d.id = @id;