package orders

import (
    "context"
    "regexp"
    "testing"

    storepkg "L0-arch/internal/db"

    sqlmockpkg "github.com/DATA-DOG/go-sqlmock"
)

func TestRepository_SavePayment(t *testing.T) {
    db, mock, err := sqlmockpkg.New()
    if err != nil { t.Fatalf("sqlmock.New: %v", err) }
    defer db.Close()

    st := storepkg.New()
    st.SetConn(db)
    repo := NewRepository(st)

    ctx := context.Background()
    orderUID := "uid-1"
    payment := Payment{
        Transaction:  "tx-1",
        RequestID:    "",
        Currency:     "RUB",
        Provider:     "wbpay",
        Amount:       1000,
        PaymentDt:    1605040200,
        Bank:         "alpha",
        DeliveryCost: 300,
        GoodsTotal:   700,
        CustomFee:    0,
    }

    mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO payment (
			order_uid, transaction_id, request_id, currency, provider,
			amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT (order_uid) DO NOTHING`)).
        WithArgs(orderUID, payment.Transaction, payment.RequestID, payment.Currency, payment.Provider, payment.Amount, payment.PaymentDt, payment.Bank, payment.DeliveryCost, payment.GoodsTotal, payment.CustomFee).
        WillReturnResult(sqlmockpkg.NewResult(0, 1))

    if err := repo.SavePayment(ctx, payment, orderUID); err != nil {
        t.Fatalf("SavePayment error: %v", err)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatalf("unmet expectations: %v", err)
    }
}


