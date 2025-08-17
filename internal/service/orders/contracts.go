package orders

import (
	"context"

	"L0-arch/internal/repository/orders"
)

type Repository interface {
	SaveOrder(ctx context.Context, o orders.Model) error
	SavePayment(ctx context.Context, p orders.Payment, orderUID string) error
	SaveDelivery(ctx context.Context, d orders.Delivery, orderUID string) error
	SaveItems(ctx context.Context, items []orders.Item, orderUID string) error
	Get(ctx context.Context, orderUID string) (*orders.Model, *orders.Delivery, *orders.Payment, *orders.Item, error)
}

type Transactor interface {
	Do(context.Context, func(context.Context) error) error
}
