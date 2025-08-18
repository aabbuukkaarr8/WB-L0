package orders

import (
	"context"

	"L0-arch/internal/db"

	trmsql "github.com/avito-tech/go-transaction-manager/drivers/sql/v2"
)

type Repository struct {
	db     *db.Store
	getter *trmsql.CtxGetter
}

func NewRepository(db *db.Store) *Repository {
	return &Repository{
		db:     db,
		getter: trmsql.DefaultCtxGetter,
	}
}

func (r *Repository) GetTr(ctx context.Context) trmsql.Tr {
	return r.getter.DefaultTrOrDB(ctx, r.db.GetConn())
}
