package orders

import (
	"sync"
)

type Service struct {
	repository Repository
	transactor Transactor
	cache      map[string]Model
	mu         sync.RWMutex
}

func NewService(
	repository Repository,
	transactor Transactor,
) *Service {
	return &Service{
		repository: repository,
		transactor: transactor,
		cache:      make(map[string]Model),
	}
}
