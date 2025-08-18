package orders

import (
	"context"
)

func (s *Service) Get(ctx context.Context, orderUID string) (*Model, error) {
	if cached, ok := s.cacheGet(orderUID); ok {
		cp := cached
		return &cp, nil
	}

	dbm, dbmD, dbmP, dbmI, err := s.repository.Get(ctx, orderUID)
	if err != nil {
		return nil, err
	}

	var m Model
	m.FillFromDB(dbm, dbmD, dbmP, dbmI)

	s.cachePut(orderUID, m)

	return &m, nil
}

 
