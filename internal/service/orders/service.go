package orders

import (
	"github.com/jmoiron/sqlx"
	"sync"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

type Service struct {
	db         *sqlx.DB
	repository Repository
	trManager  *manager.Manager
	cache      map[string]Model
	mu         sync.RWMutex
}

func NewService(
	repository Repository,
	dbConn *sqlx.DB,
) *Service {
	return &Service{
		repository: repository,
		cache:      make(map[string]Model),
		trManager:  manager.Must(trmsqlx.NewDefaultFactory(dbConn)),
	}
}
