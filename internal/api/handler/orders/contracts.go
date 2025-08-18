package orders

import (
	"L0-arch/internal/service/orders"
	"context"
)

type Service interface {
	Get(ctx context.Context, orderUID string) (*orders.Model, error)
	SaveOrder(ctx context.Context, o orders.Model) error
}
