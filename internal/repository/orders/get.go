package orders

func (r *RepositoryGet) Get(orderUID string) (*Model, *Delivery, *Payment, []Item, error) {
	o := &Model{}
	d := &Delivery{}
	p := &Payment{}
	items := make([]Item, 0)

	const qOrder = `
		SELECT
			order_uid, track_number, entry, locale, internal_signature,
			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders
		WHERE order_uid = $1
	`
	if err := r.db.GetConn().QueryRow(qOrder, orderUID).Scan(
		&o.OrderUID,
		&o.TrackNumber,
		&o.Entry,
		&o.Locale,
		&o.InternalSignature,
		&o.CustomerID,
		&o.DeliveryService,
		&o.ShardKey,
		&o.SmID,
		&o.DateCreated,
		&o.OofShard,
	); err != nil {
		return nil, nil, nil, nil, err
	}

	// delivery
	const qDelivery = `
		SELECT name, phone, zip, city, address, region, email
		FROM delivery
		WHERE order_uid = $1
	`
	if err := r.db.GetConn().QueryRow(qDelivery, orderUID).Scan(
		&d.Name,
		&d.Phone,
		&d.Zip,
		&d.City,
		&d.Address,
		&d.Region,
		&d.Email,
	); err != nil {
		return nil, nil, nil, nil, err
	}

	// payment
	const qPayment = `
		SELECT transaction_id, request_id, currency, provider,
		       amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		FROM payment
		WHERE order_uid = $1
	`
	if err := r.db.GetConn().QueryRow(qPayment, orderUID).Scan(
		&p.Transaction,
		&p.RequestID,
		&p.Currency,
		&p.Provider,
		&p.Amount,
		&p.PaymentDt,
		&p.Bank,
		&p.DeliveryCost,
		&p.GoodsTotal,
		&p.CustomFee,
	); err != nil {
		return nil, nil, nil, nil, err
	}

	const qItems = `
		SELECT chrt_id, track_number, price, rid, name,
		       sale, size, total_price, nm_id, brand, status
		FROM items
		WHERE order_uid = $1
	`
	rows, err := r.db.GetConn().Query(qItems, orderUID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var it Item
		if err := rows.Scan(
			&it.ChrtID,
			&it.TrackNumber,
			&it.Price,
			&it.RID,
			&it.Name,
			&it.Sale,
			&it.Size,
			&it.TotalPrice,
			&it.NmID,
			&it.Brand,
			&it.Status,
		); err != nil {
			return nil, nil, nil, nil, err
		}
		items = append(items, it)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, nil, nil, err
	}

	return o, d, p, items, nil
}
