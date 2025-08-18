CREATE TABLE items (
                       chrt_id      BIGINT PRIMARY KEY,
                       order_uid    TEXT NOT NULL
                           REFERENCES orders(order_uid)
                               ON DELETE CASCADE,
                       track_number TEXT NOT NULL,
                       price        INTEGER NOT NULL,
                       rid          TEXT NOT NULL,
                       name         TEXT NOT NULL,
                       sale         INTEGER NOT NULL,
                       size         TEXT NOT NULL,
                       total_price  INTEGER NOT NULL,
                       nm_id        INTEGER NOT NULL,
                       brand        TEXT NOT NULL,
                       status       INTEGER NOT NULL
);

