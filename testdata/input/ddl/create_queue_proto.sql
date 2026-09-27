CREATE QUEUE Q (id INT64, Payload examples.Message NOT NULL) PRIMARY KEY (id), OPTIONS (disable_delivery = true)
