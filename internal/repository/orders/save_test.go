package orders

import (
    "context"
    "regexp"
    "testing"

    storepkg "L0-arch/internal/db"

    sqlmockpkg "github.com/DATA-DOG/go-sqlmock"
)

func TestRepository_SaveOrder(t *testing.T) {
    db, mock, err := sqlmockpkg.New()
    if err != nil { t.Fatalf("sqlmock.New: %v", err) }
    defer db.Close()

    st := storepkg.New()
    st.SetConn(db)
    repo := NewRepository(st)

    ctx := context.Background()
    order := Model{
        OrderUID:          "uid-1",
        TrackNumber:       "TR-1",
        Entry:             "WBIL",
        Locale:            "ru",
        InternalSignature: "",
        CustomerID:        "cust-1",
        DeliveryService:   "meest",
        ShardKey:          "9",
        SmID:              99,
        DateCreated:       "2020-11-10T21:30:00Z",
        OofShard:          "1",
    }

    mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO orders (
			order_uid, track_number, entry, locale, internal_signature,
			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT (order_uid) DO NOTHING`)).
        WithArgs(
            order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
            order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard,
        ).
        WillReturnResult(sqlmockpkg.NewResult(0, 1))

    if err := repo.SaveOrder(ctx, order); err != nil {
        t.Fatalf("SaveOrder error: %v", err)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatalf("unmet expectations: %v", err)
    }
}


