CREATE TABLE IF NOT EXISTS user_goals (
    user_id INT PRIMARY KEY,
    calories INT NOT NULL,
    proteins INT NOT NULL,
    fats INT NOT NULL,
    carbohydrates INT NOT NULL
);
