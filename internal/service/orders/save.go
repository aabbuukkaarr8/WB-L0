package orders

import (
	"context"
	"fmt"
	"time"

	"L0-arch/pkg/retry"
)

func (s *Service) SaveOrder(ctx context.Context, o Model) error {
	toDB := s.Converter(o, o.Delivery, o.Payment, o.Items)

	op := func() error {
		return s.transactor.Do(ctx, func(ctx context.Context) error {
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
	}
	err := retry.Do(ctx, op, retry.Options{
		MaxAttempts: 5,
		BaseDelay:   100 * time.Millisecond,
		MaxDelay:    2 * time.Second,
		IsRetryable: retry.IsPgRetryable,
	})
	if err != nil {
		return fmt.Errorf("save order: %w", err)
	}

	s.cachePut(o.OrderUID, o)
	return nil
}
