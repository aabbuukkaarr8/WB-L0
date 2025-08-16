CREATE TABLE delivery (
                          order_uid TEXT PRIMARY KEY
                              REFERENCES orders(order_uid)
                                  ON DELETE CASCADE,
                          name    TEXT NOT NULL,
                          phone   TEXT NOT NULL,
                          zip     TEXT NOT NULL,
                          city    TEXT NOT NULL,
                          address TEXT NOT NULL,
                          region  TEXT NOT NULL,
                          email   TEXT NOT NULL
);
