package orders

import "L0-arch/internal/service/orders"

type Service interface {
	Get(orderUID string) (*orders.Model, error)
}
