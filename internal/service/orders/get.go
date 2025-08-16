package orders

func (s *Service) Get(orderUID string) (*Model, error) {
	s.mu.RLock()
	if cached, ok := s.cache[orderUID]; ok {
		s.mu.RUnlock()
		cp := cached
		return &cp, nil
	}
	s.mu.RUnlock()

	dbm, dbmD, dbmP, dbmI, err := s.repository.GetOrder(orderUID)
	if err != nil {
		return nil, err
	}

	var m Model
	m.FillFromDB(dbm, dbmD, dbmP, dbmI)

	s.mu.Lock()
	s.cache[orderUID] = m
	s.mu.Unlock()

	return &m, nil
}

//func (s *Service) Get(orderUID string) (*Model, error) {
//	dbm, dbmD, dbmP, dbmI, err := s.repository.Get(orderUID)
//	if err != nil {
//		return nil, err
//	}
//
//	p := &Model{}
//	var items []orders.Item
//	if dbmI != nil {
//		items = []orders.Item{*dbmI}
//	}
//	p.FillFromDB(dbm, dbmD, dbmP, &items)
//	return p, nil
//}

//func (s *Service) Get(orderUID string) (*Model, error) {
//	p := &Model{}
//
//	r := orders.Model{}
//	r = s.cache[orderUID]
//	p.FillFromDB(&r, &r.Delivery, &r.Payment, &r.Items)
//	if p == nil {
//		dbm, dbmD, dbmP, dbmI, err := s.repository.Get(orderUID)
//		if err != nil {
//			return nil, err
//		}
//		var items []orders.Item
//		if dbmI != nil {
//			items = []orders.Item{*dbmI}
//		}
//		p.FillFromDB(dbm, dbmD, dbmP, &items)
//	}
//	return p, nil
//}
