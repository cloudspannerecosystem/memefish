CREATE QUEUE queue (constraint INT64, foreign INT64, synonym INT64, Payload BYTES(MAX) NOT NULL) PRIMARY KEY (constraint, foreign, synonym)
