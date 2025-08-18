package orders

import "L0-arch/internal/repository/orders"

type Model struct {
	OrderUID          string   `json:"order_uid" validate:"required"`
	TrackNumber       string   `json:"track_number" validate:"required"`
	Entry             string   `json:"entry" validate:"required"`
	Delivery          Delivery `json:"delivery" validate:"required"`
	Payment           Payment  `json:"payment" validate:"required"`
	Items             []Item   `json:"items" validate:"min=1,dive"`
	Locale            string   `json:"locale" validate:"required"`
	InternalSignature string   `json:"internal_signature"`
	CustomerID        string   `json:"customer_id" validate:"required"`
	DeliveryService   string   `json:"delivery_service" validate:"required"`
	ShardKey          string   `json:"shardkey" validate:"required"`
	SmID              int      `json:"sm_id" validate:"gte=0"`
	DateCreated       string   `json:"date_created" validate:"required"`
	OofShard          string   `json:"oof_shard" validate:"required"`
}

type Delivery struct {
	Name    string `json:"name" validate:"not_empty"`
	Phone   string `json:"phone" validate:"not_empty"`
	Zip     string `json:"zip" validate:"not_empty"`
	City    string `json:"city" validate:"not_empty"`
	Address string `json:"address" validate:"not_empty"`
	Region  string `json:"region" validate:"not_empty"`
	Email   string `json:"email" validate:"required,email"`
}

type Payment struct {
	Transaction  string `json:"transaction" validate:"not_empty"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency" validate:"not_empty"`
	Provider     string `json:"provider" validate:"not_empty"`
	Amount       int    `json:"amount" validate:"gte=0"`
	PaymentDt    int64  `json:"payment_dt" validate:"gte=0"`
	Bank         string `json:"bank" validate:"not_empty"`
	DeliveryCost int    `json:"delivery_cost" validate:"gte=0"`
	GoodsTotal   int    `json:"goods_total" validate:"gte=0"`
	CustomFee    int    `json:"custom_fee" validate:"gte=0"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id" validate:"gt=0"`
	TrackNumber string `json:"track_number" validate:"not_empty"`
	Price       int    `json:"price" validate:"gte=0"`
	RID         string `json:"rid" validate:"not_empty"`
	Name        string `json:"name" validate:"not_empty"`
	Sale        int    `json:"sale" validate:"gte=0"`
	Size        string `json:"size" validate:"not_empty"`
	TotalPrice  int    `json:"total_price" validate:"gte=0"`
	NmID        int    `json:"nm_id" validate:"gte=0"`
	Brand       string `json:"brand" validate:"not_empty"`
	Status      int    `json:"status" validate:"gte=0"`
}

func (m *Model) FillFromDB(dbm *orders.Model, dbd *orders.Delivery, dbp *orders.Payment, dbi []orders.Item) {
	m.Delivery = Delivery{
		Name:    dbd.Name,
		Phone:   dbd.Phone,
		Zip:     dbd.Zip,
		City:    dbd.City,
		Address: dbd.Address,
		Region:  dbd.Region,
		Email:   dbd.Email,
	}
	m.Payment = Payment{
		Transaction:  dbp.Transaction,
		RequestID:    dbp.RequestID,
		Currency:     dbp.Currency,
		Provider:     dbp.Provider,
		Amount:       dbp.Amount,
		PaymentDt:    dbp.PaymentDt,
		Bank:         dbp.Bank,
		DeliveryCost: dbp.DeliveryCost,
		GoodsTotal:   dbp.GoodsTotal,
		CustomFee:    dbp.CustomFee,
	}
	m.Items = m.Items[:0]
	if dbi != nil {
		m.Items = make([]Item, 0, len(dbi))
		for _, it := range dbi {
			m.Items = append(m.Items, Item{
				ChrtID:      it.ChrtID,
				TrackNumber: it.TrackNumber,
				Price:       it.Price,
				RID:         it.RID,
				Name:        it.Name,
				Sale:        it.Sale,
				Size:        it.Size,
				TotalPrice:  it.TotalPrice,
				NmID:        it.NmID,
				Brand:       it.Brand,
				Status:      it.Status,
			})
		}
	}
	m.OrderUID = dbm.OrderUID
	m.TrackNumber = dbm.TrackNumber
	m.Entry = dbm.Entry
	m.Locale = dbm.Locale
	m.InternalSignature = dbm.InternalSignature
	m.CustomerID = dbm.CustomerID
	m.DeliveryService = dbm.DeliveryService
	m.ShardKey = dbm.ShardKey
	m.SmID = dbm.SmID
	m.DateCreated = dbm.DateCreated
	m.OofShard = dbm.OofShard

}
func (s *Service) Converter(o Model, d Delivery, p Payment, items []Item) orders.Model {
	SVCDelivery := orders.Delivery{
		Name:    d.Name,
		Phone:   d.Phone,
		Zip:     d.Zip,
		City:    d.City,
		Address: d.Address,
		Region:  d.Region,
		Email:   d.Email,
	}
	SVCPayment := orders.Payment{
		Transaction:  p.Transaction,
		RequestID:    p.RequestID,
		Currency:     p.Currency,
		Provider:     p.Provider,
		Amount:       p.Amount,
		PaymentDt:    p.PaymentDt,
		Bank:         p.Bank,
		DeliveryCost: p.DeliveryCost,
		GoodsTotal:   p.GoodsTotal,
		CustomFee:    p.CustomFee,
	}
	SVCItems := make([]orders.Item, 0, len(items))
	for _, it := range items {
		SVCItems = append(SVCItems, orders.Item{
			ChrtID:      it.ChrtID,
			TrackNumber: it.TrackNumber,
			Price:       it.Price,
			RID:         it.RID,
			Name:        it.Name,
			Sale:        it.Sale,
			Size:        it.Size,
			TotalPrice:  it.TotalPrice,
			NmID:        it.NmID,
			Brand:       it.Brand,
			Status:      it.Status,
		})
	}

	SVCOrder := orders.Model{
		OrderUID:          o.OrderUID,
		TrackNumber:       o.TrackNumber,
		Entry:             o.Entry,
		Delivery:          SVCDelivery,
		Payment:           SVCPayment,
		Items:             SVCItems,
		Locale:            o.Locale,
		InternalSignature: o.InternalSignature,
		CustomerID:        o.CustomerID,
		DeliveryService:   o.DeliveryService,
		ShardKey:          o.ShardKey,
		SmID:              o.SmID,
		DateCreated:       o.DateCreated,
		OofShard:          o.OofShard,
	}

	return SVCOrder

}
