package main

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	repo "github.com/vacmannnn/calorina/internal/adapter/sqlite/dishes"
)

func initDB() *repo.Repository {
	db, err := sql.Open("sqlite3", "file:mydb.db?cache=shared&mode=rwc")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Apply migrations
	m, err := migrate.New(
		"file://schema/sqlite", // Path to your migration files
		"sqlite3://mydb.db")    // Database URL
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully!")

	return repo.NewRepository(db)
}
