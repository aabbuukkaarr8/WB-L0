package orders

import (
	"L0-arch/internal/db"
)

type Repository struct {
	db *db.Store
}

func NewRepository(db *db.Store) *Repository {
	return &Repository{
		db: db,
	}
}
