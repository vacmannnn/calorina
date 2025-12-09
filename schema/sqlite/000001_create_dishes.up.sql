CREATE TABLE IF NOT EXISTS dishes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    calories INT NOT NULL,
    protein INT DEFAULT '0',
    fat INT DEFAULT '0',
    carbohydrates INT DEFAULT '0',
    weight INT DEFAULT '0'
);

CREATE TABLE IF NOT EXISTS daily_eating_info (
    user_id INT NOT NULL,
    date TEXT NOT NULL,
    total_calories INT NOT NULL,
    total_proteins INT,
    total_fats INT,
    total_carbohydrates INT,
    dishes_ids TEXT NOT NULL,
    PRIMARY KEY (user_id, date)
);