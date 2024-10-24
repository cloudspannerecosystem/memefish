CREATE TABLE Singers (
    SingerId INT64 NOT NULL,
    SingerName STRING(1024),
    SYNONYM (Artists)
) PRIMARY KEY (SingerId)