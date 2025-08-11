package orders

import (
	"sync"
)

type Service struct {
	repository Repository
	cache      map[string]Model
	mu         sync.RWMutex
}

func NewService(
	repository Repository,
) *Service {
	return &Service{
		repository: repository,
		cache:      make(map[string]Model),
	}
}
