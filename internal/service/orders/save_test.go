package orders

import (
	"context"
	"errors"
	"testing"

	repo "L0-arch/internal/repository/orders"
)

type fakeRepoSave struct {
	calls map[string]int
}

func sampleServiceModel() Model {
	return Model{
		OrderUID:    "uid-1",
		TrackNumber: "TR-1",
		Entry:       "WBIL",
		Delivery: Delivery{
			Name:    "John Doe",
			Phone:   "+79990000000",
			Zip:     "101000",
			City:    "Moscow",
			Address: "Tverskaya 1",
			Region:  "Moscow",
			Email:   "john@example.com",
		},
		Payment: Payment{
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
		},
		Items: []Item{{
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
		}},
		Locale:            "ru",
		InternalSignature: "",
		CustomerID:        "cust-1",
		DeliveryService:   "meest",
		ShardKey:          "9",
		SmID:              99,
		DateCreated:       "2020-11-10T21:30:00Z",
		OofShard:          "1",
	}
}

func (f *fakeRepoSave) SaveOrder(ctx context.Context, o repo.Model) error {
	f.calls["SaveOrder"]++
	return nil
}
func (f *fakeRepoSave) SavePayment(ctx context.Context, p repo.Payment, orderUID string) error {
	f.calls["SavePayment"]++
	return nil
}
func (f *fakeRepoSave) SaveDelivery(ctx context.Context, d repo.Delivery, orderUID string) error {
	f.calls["SaveDelivery"]++
	return nil
}
func (f *fakeRepoSave) SaveItems(ctx context.Context, items []repo.Item, orderUID string) error {
	f.calls["SaveItems"]++
	return nil
}
func (f *fakeRepoSave) Get(ctx context.Context, orderUID string) (*repo.Model, *repo.Delivery, *repo.Payment, []repo.Item, error) {
	return nil, nil, nil, nil, errors.New("not impl")
}

type fakeTxSave struct{ ret error }

func (t *fakeTxSave) Do(ctx context.Context, fn func(context.Context) error) error {
	if t.ret != nil {
		return t.ret
	}
	return fn(ctx)
}

func TestService_SaveOrder(t *testing.T) {
	svc := NewService(&fakeRepoSave{calls: map[string]int{}}, &fakeTxSave{})
	if err := svc.SaveOrder(context.Background(), sampleServiceModel()); err != nil {
		t.Fatalf("SaveOrder error: %v", err)
	}
}

func TestService_SaveOrder_TxErr(t *testing.T) {
	svc := NewService(&fakeRepoSave{calls: map[string]int{}}, &fakeTxSave{ret: errors.New("tx fail")})
	if err := svc.SaveOrder(context.Background(), sampleServiceModel()); err == nil {
		t.Fatalf("expected error, got nil")
	}
}
