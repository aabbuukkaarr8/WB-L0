package orders

import (
	"L0-arch/internal/db"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

type RepositoryGet struct {
	db *db.Store
}

func NewRepositoryGet(db *db.Store) *RepositoryGet {
	return &RepositoryGet{
		db: db,
	}
}
