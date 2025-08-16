CREATE TABLE payment (
                         order_uid      TEXT PRIMARY KEY
                             REFERENCES orders(order_uid)
                                 ON DELETE CASCADE,
                         transaction_id TEXT NOT NULL,
                         request_id     TEXT,
                         currency       TEXT NOT NULL,
                         provider       TEXT NOT NULL,
                         amount         INTEGER NOT NULL,
                         payment_dt     BIGINT NOT NULL,
                         bank           TEXT NOT NULL,
                         delivery_cost  INTEGER NOT NULL,
                         goods_total    INTEGER NOT NULL,
                         custom_fee     INTEGER NOT NULL
);
