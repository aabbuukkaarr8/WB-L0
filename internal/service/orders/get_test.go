package orders

import (
    "context"
    "errors"
    "reflect"
    "testing"

    repo "L0-arch/internal/repository/orders"
)

type fakeRepoGet struct{
    o *repo.Model
    d *repo.Delivery
    p *repo.Payment
    items []repo.Item
    err error
}
func (f *fakeRepoGet) SaveOrder(ctx context.Context, o repo.Model) error { return nil }
func (f *fakeRepoGet) SavePayment(ctx context.Context, p repo.Payment, orderUID string) error { return nil }
func (f *fakeRepoGet) SaveDelivery(ctx context.Context, d repo.Delivery, orderUID string) error { return nil }
func (f *fakeRepoGet) SaveItems(ctx context.Context, items []repo.Item, orderUID string) error { return nil }
func (f *fakeRepoGet) Get(ctx context.Context, orderUID string) (*repo.Model, *repo.Delivery, *repo.Payment, []repo.Item, error) { return f.o, f.d, f.p, f.items, f.err }

type fakeTxGet struct{}
func (t *fakeTxGet) Do(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }

func TestService_Get_Cached(t *testing.T) {
    svc := NewService(&fakeRepoGet{err: errors.New("should not be called")}, &fakeTxGet{})
    m := sampleServiceModel()
    svc.cache[m.OrderUID] = m
    got, err := svc.Get(context.Background(), m.OrderUID)
    if err != nil { t.Fatalf("Get error: %v", err) }
    if !reflect.DeepEqual(*got, m) { t.Fatalf("unexpected cached value") }
}

func TestService_Get_DB(t *testing.T) {
    m := sampleServiceModel()
    fr := &fakeRepoGet{
        o: &repo.Model{
            OrderUID: m.OrderUID,
            TrackNumber: m.TrackNumber,
            Entry: m.Entry,
            Locale: m.Locale,
            InternalSignature: m.InternalSignature,
            CustomerID: m.CustomerID,
            DeliveryService: m.DeliveryService,
            ShardKey: m.ShardKey,
            SmID: m.SmID,
            DateCreated: m.DateCreated,
            OofShard: m.OofShard,
        },
        d: &repo.Delivery{
            Name: m.Delivery.Name,
            Phone: m.Delivery.Phone,
            Zip: m.Delivery.Zip,
            City: m.Delivery.City,
            Address: m.Delivery.Address,
            Region: m.Delivery.Region,
            Email: m.Delivery.Email,
        },
        p: &repo.Payment{
            Transaction: m.Payment.Transaction,
            RequestID: m.Payment.RequestID,
            Currency: m.Payment.Currency,
            Provider: m.Payment.Provider,
            Amount: m.Payment.Amount,
            PaymentDt: m.Payment.PaymentDt,
            Bank: m.Payment.Bank,
            DeliveryCost: m.Payment.DeliveryCost,
            GoodsTotal: m.Payment.GoodsTotal,
            CustomFee: m.Payment.CustomFee,
        },
        items: []repo.Item{{
            ChrtID: m.Items[0].ChrtID,
            TrackNumber: m.Items[0].TrackNumber,
            Price: m.Items[0].Price,
            RID: m.Items[0].RID,
            Name: m.Items[0].Name,
            Sale: m.Items[0].Sale,
            Size: m.Items[0].Size,
            TotalPrice: m.Items[0].TotalPrice,
            NmID: m.Items[0].NmID,
            Brand: m.Items[0].Brand,
            Status: m.Items[0].Status,
        }},
    }
    svc := NewService(fr, &fakeTxGet{})
    got, err := svc.Get(context.Background(), m.OrderUID)
    if err != nil { t.Fatalf("Get error: %v", err) }
    if got.OrderUID != m.OrderUID || len(got.Items) != 1 {
        t.Fatalf("unexpected result: %#v", got)
    }
}


