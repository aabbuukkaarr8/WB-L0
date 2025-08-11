package orders

import (
	"L0-arch/internal/repository/orders"
)

type Repository interface {
	SaveOrder(o orders.Model) error
	SavePayment(p orders.Payment, orderUID string) error
	SaveDelivery(d orders.Delivery, orderUID string) error
	SaveItems(items []orders.Item, orderUID string) error
	Get(orderUID string) (*orders.Model, *orders.Delivery, *orders.Payment, *orders.Item, error)
}
