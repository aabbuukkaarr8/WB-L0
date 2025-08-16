package orders

import (
	"L0-arch/internal/repository/orders"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	SaveOrderTx(tx *sqlx.Tx, o orders.Model) error
	SavePaymentTx(tx *sqlx.Tx, p orders.Payment, orderUID string) error
	SaveDeliveryTx(tx *sqlx.Tx, d orders.Delivery, orderUID string) error
	SaveItemsTx(tx *sqlx.Tx, items []orders.Item, orderUID string) error
	GetOrder(orderUID string) (*orders.Model, *orders.Delivery, *orders.Payment, []orders.Item, error)
}
