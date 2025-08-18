package orders

import "L0-arch/internal/service/orders"

type Model struct {
	OrderUID          string   `json:"order_uid"`
	TrackNumber       string   `json:"track_number"`
	Entry             string   `json:"entry"`
	Delivery          Delivery `json:"delivery"`
	Payment           Payment  `json:"payment"`
	Items             []Item   `json:"items"`
	Locale            string   `json:"locale"`
	InternalSignature string   `json:"internal_signature"`
	CustomerID        string   `json:"customer_id"`
	DeliveryService   string   `json:"delivery_service"`
	ShardKey          string   `json:"shardkey"`
	SmID              int      `json:"sm_id"`
	DateCreated       string   `json:"date_created"`
	OofShard          string   `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func (m *Model) FillFromService(dbm *orders.Model, dbd *orders.Delivery, dbp *orders.Payment, dbi *[]orders.Item) {
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
		m.Items = make([]Item, 0, len(*dbi))
		for _, it := range *dbi {
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
