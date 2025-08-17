package orders

import "context"

func (r *Repository) SaveItems(ctx context.Context, items []Item, orderUID string) error {
	for _, it := range items {
		if _, err := r.GetTr(ctx).Exec(`
			INSERT INTO items (
				chrt_id, order_uid, track_number, price, rid, name,
				sale, size, total_price, nm_id, brand, status
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)	
			ON CONFLICT (chrt_id) DO NOTHING
		`,
			it.ChrtID,
			orderUID,
			it.TrackNumber,
			it.Price,
			it.RID,
			it.Name,
			it.Sale,
			it.Size,
			it.TotalPrice,
			it.NmID,
			it.Brand,
			it.Status,
		); err != nil {
			return err
		}
	}
	return nil
}
