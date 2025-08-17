package orders

import (
	"context"
	"fmt"
)

func (s *Service) SaveOrder(ctx context.Context, o Model) error {
	toDB := s.Converter(o, o.Delivery, o.Payment, o.Items)

	err := s.transactor.Do(ctx, func(ctx context.Context) error {
		err := s.repository.SaveOrder(ctx, toDB)
		if err != nil {
			return err
		}
		err = s.repository.SaveItems(ctx, toDB.Items, toDB.OrderUID)
		if err != nil {
			return err
		}
		err = s.repository.SaveDelivery(ctx, toDB.Delivery, toDB.OrderUID)
		if err != nil {
			return err
		}
		err = s.repository.SavePayment(ctx, toDB.Payment, toDB.OrderUID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("save order: %w", err)
	}

	s.mu.Lock()
	s.cache[o.OrderUID] = o
	s.mu.Unlock()
	return nil
}
