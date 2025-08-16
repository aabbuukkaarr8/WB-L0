package orders

import (
	"context"
)

func (s *Service) SaveOrder(ctx context.Context, o Model) error {
	toDB := s.Converter(o, o.Delivery, o.Payment, o.Items)
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = s.repository.SaveOrderTx(tx, toDB); err != nil {
		return err
	}
	if err = s.repository.SaveItemsTx(tx, toDB.Items, toDB.OrderUID); err != nil {
		return err
	}
	if err = s.repository.SaveDeliveryTx(tx, toDB.Delivery, toDB.OrderUID); err != nil {
		return err
	}
	if err = s.repository.SavePaymentTx(tx, toDB.Payment, toDB.OrderUID); err != nil {
		return err
	}
	s.mu.Lock()
	s.cache[o.OrderUID] = o
	s.mu.Unlock()
	return nil
}
