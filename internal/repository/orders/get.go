package orders

func (r *Repository) Get(orderUID string) (*Model, *Delivery, *Payment, *Item, error) {
	o := &Model{}
	d := &Delivery{}
	p := &Payment{}
	i := &Item{}
	queryM := `SELECT * FROM orders WHERE order_uid = $1`
	err := r.db.GetConn().QueryRow(queryM, orderUID).Scan(o)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	queryD := `SELECT * FROM delivery WHERE order_uid = $1`
	err = r.db.GetConn().QueryRow(queryD, orderUID).Scan(d)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	queryI := `SELECT * FROM items WHERE order_uid = $1`
	err = r.db.GetConn().QueryRow(queryI, orderUID).Scan(i)
	if err != nil {
		return nil, nil, nil, nil, err

	}
	queryP := `SELECT * FROM payment WHERE order_uid = $1`
	err = r.db.GetConn().QueryRow(queryP, orderUID).Scan(p)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return o, d, p, i, nil

}
