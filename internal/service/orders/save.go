package orders

func (s *Service) SaveOrder(o Model) error {
	toDB := s.Converter(o, o.Delivery, o.Payment, o.Items)

	err := s.repository.SaveOrder(toDB)
	if err != nil {
		return err
	}
	err = s.repository.SaveItems(toDB.Items, toDB.OrderUID)
	if err != nil {
		return err
	}
	err = s.repository.SaveDelivery(toDB.Delivery, toDB.OrderUID)
	if err != nil {
		return err
	}
	err = s.repository.SavePayment(toDB.Payment, toDB.OrderUID)
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.cache[o.OrderUID] = o
	s.mu.Unlock()
	return nil
}
