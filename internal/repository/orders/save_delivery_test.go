package orders

import (
    "context"
    "regexp"
    "testing"

    storepkg "L0-arch/internal/db"

    sqlmockpkg "github.com/DATA-DOG/go-sqlmock"
)

func TestRepository_SaveDelivery(t *testing.T) {
    db, mock, err := sqlmockpkg.New()
    if err != nil { t.Fatalf("sqlmock.New: %v", err) }
    defer db.Close()

    st := storepkg.New()
    st.SetConn(db)
    repo := NewRepository(st)

    ctx := context.Background()
    orderUID := "uid-1"
    delivery := Delivery{
        Name:    "John Doe",
        Phone:   "+79990000000",
        Zip:     "101000",
        City:    "Moscow",
        Address: "Tverskaya 1",
        Region:  "Moscow",
        Email:   "john@example.com",
    }

    mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO delivery (
			order_uid, name, phone, zip, city, address, region, email
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		ON CONFLICT (order_uid) DO NOTHING`)).
        WithArgs(orderUID, delivery.Name, delivery.Phone, delivery.Zip, delivery.City, delivery.Address, delivery.Region, delivery.Email).
        WillReturnResult(sqlmockpkg.NewResult(0, 1))

    if err := repo.SaveDelivery(ctx, delivery, orderUID); err != nil {
        t.Fatalf("SaveDelivery error: %v", err)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatalf("unmet expectations: %v", err)
    }
}


