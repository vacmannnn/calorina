package dishes

import (
	"database/sql"

	"github.com/vacmannnn/calorina/internal/adapter/sqlite/dishes/sqlc"
)

type Repository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db:      db,
		queries: sqlc.New(db),
	}
}
