package orders

func (r *Repository) SaveOrder(o Model) error {
	_, err := r.db.GetConn().Exec(`
		INSERT INTO orders (
			order_uid, track_number, entry, locale, internal_signature,
			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT (order_uid) DO NOTHING
	`,
		o.OrderUID, o.TrackNumber, o.Entry, o.Locale, o.InternalSignature,
		o.CustomerID, o.DeliveryService, o.ShardKey, o.SmID, o.DateCreated, o.OofShard,
	)
	if err != nil {
		return err
	}
	return nil

}
