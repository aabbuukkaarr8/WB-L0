package orders

import (
    "context"
    "regexp"
    "testing"

    storepkg "L0-arch/internal/db"

    sqlmockpkg "github.com/DATA-DOG/go-sqlmock"
)

func TestRepository_SaveItems(t *testing.T) {
    db, mock, err := sqlmockpkg.New()
    if err != nil { t.Fatalf("sqlmock.New: %v", err) }
    defer db.Close()

    st := storepkg.New()
    st.SetConn(db)
    repo := NewRepository(st)

    ctx := context.Background()
    orderUID := "uid-1"
    items := []Item{{
        ChrtID:      1,
        TrackNumber: "TR-1",
        Price:       700,
        RID:         "rid-1",
        Name:        "T-Shirt",
        Sale:        0,
        Size:        "M",
        TotalPrice:  700,
        NmID:        123,
        Brand:       "WB",
        Status:      1,
    }}

    mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO items (
				chrt_id, order_uid, track_number, price, rid, name,
				sale, size, total_price, nm_id, brand, status
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)	
			ON CONFLICT (chrt_id) DO NOTHING`)).
        WithArgs(
            items[0].ChrtID,
            orderUID,
            items[0].TrackNumber,
            items[0].Price,
            items[0].RID,
            items[0].Name,
            items[0].Sale,
            items[0].Size,
            items[0].TotalPrice,
            items[0].NmID,
            items[0].Brand,
            items[0].Status,
        ).
        WillReturnResult(sqlmockpkg.NewResult(0, 1))

    if err := repo.SaveItems(ctx, items, orderUID); err != nil {
        t.Fatalf("SaveItems error: %v", err)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatalf("unmet expectations: %v", err)
    }
}


