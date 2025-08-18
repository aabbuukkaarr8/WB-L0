package orders

import (
    "context"
    "regexp"
    "testing"

    storepkg "L0-arch/internal/db"

    sqlmockpkg "github.com/DATA-DOG/go-sqlmock"
)

func TestRepository_Get(t *testing.T) {
    db, mock, err := sqlmockpkg.New()
    if err != nil { t.Fatalf("sqlmock.New: %v", err) }
    defer db.Close()

    st := storepkg.New()
    st.SetConn(db)
    repo := NewRepository(st)

    ctx := context.Background()
    orderUID := "uid-1"

    rowsOrder := sqlmockpkg.NewRows([]string{
        "order_uid", "track_number", "entry", "locale", "internal_signature",
        "customer_id", "delivery_service", "shardkey", "sm_id", "date_created", "oof_shard",
    }).AddRow(
        orderUID, "TR-1", "WBIL", "ru", "", "cust-1", "meest", "9", 99, "2020-11-10T21:30:00Z", "1",
    )
    mock.ExpectQuery(regexp.QuoteMeta(`SELECT
		order_uid, track_number, entry, locale, internal_signature,
		customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
	FROM orders
	WHERE order_uid = $1
	`)).WithArgs(orderUID).WillReturnRows(rowsOrder)

    rowsDelivery := sqlmockpkg.NewRows([]string{"name", "phone", "zip", "city", "address", "region", "email"}).
        AddRow("John Doe", "+79990000000", "101000", "Moscow", "Tverskaya 1", "Moscow", "john@example.com")
    mock.ExpectQuery(regexp.QuoteMeta(`SELECT name, phone, zip, city, address, region, email
		FROM delivery
		WHERE order_uid = $1
	`)).WithArgs(orderUID).WillReturnRows(rowsDelivery)

    rowsItems := sqlmockpkg.NewRows([]string{
        "chrt_id", "track_number", "price", "rid", "name",
        "sale", "size", "total_price", "nm_id", "brand", "status",
    }).AddRow(1, "TR-1", 700, "rid-1", "T-Shirt", 0, "M", 700, 123, "WB", 1)
    mock.ExpectQuery(regexp.QuoteMeta(`SELECT chrt_id, track_number, price, rid, name,
	       sale, size, total_price, nm_id, brand, status
		FROM items
		WHERE order_uid = $1
	`)).WithArgs(orderUID).WillReturnRows(rowsItems)

    rowsPayment := sqlmockpkg.NewRows([]string{
        "transaction_id", "request_id", "currency", "provider",
        "amount", "payment_dt", "bank", "delivery_cost", "goods_total", "custom_fee",
    }).AddRow("tx-1", "", "RUB", "wbpay", 1000, 1605040200, "alpha", 300, 700, 0)
    mock.ExpectQuery(regexp.QuoteMeta(`SELECT transaction_id, request_id, currency, provider,
	       amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		FROM payment
		WHERE order_uid = $1
	`)).WithArgs(orderUID).WillReturnRows(rowsPayment)

    o, d, p, items, err := repo.Get(ctx, orderUID)
    if err != nil {
        t.Fatalf("Get error: %v", err)
    }

    if o.OrderUID != orderUID || d.Name == "" || p.Transaction == "" || len(items) != 1 {
        t.Fatalf("unexpected result: %#v %#v %#v len(items)=%d", o, d, p, len(items))
    }
}


