package orders

import "context"

func (r *Repository) SavePayment(ctx context.Context, p Payment, orderUID string) error {
	_, err := r.GetTr(ctx).Exec(`
		INSERT INTO payment (
			order_uid, transaction_id, request_id, currency, provider,
			amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT (order_uid) DO NOTHING
	`,
		orderUID,
		p.Transaction,
		p.RequestID,
		p.Currency,
		p.Provider,
		p.Amount,
		p.PaymentDt,
		p.Bank,
		p.DeliveryCost,
		p.GoodsTotal,
		p.CustomFee,
	)
	if err != nil {
		return err
	}
	return nil
}
