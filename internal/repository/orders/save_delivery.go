package orders

import "context"

func (r *Repository) SaveDelivery(ctx context.Context, d Delivery, orderUID string) error {
	_, err := r.GetTr(ctx).Exec(`
		INSERT INTO delivery (
			order_uid, name, phone, zip, city, address, region, email
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		ON CONFLICT (order_uid) DO NOTHING
	`,
		orderUID, d.Name, d.Phone, d.Zip, d.City, d.Address, d.Region, d.Email,
	)
	if err != nil {
		return err
	}
	return nil
}
