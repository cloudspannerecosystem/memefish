-- It is still valid CREATE TABLE statement.
CREATE TABLE Singers (
    SYNONYM (Ignored),
    SingerId INT64 NOT NULL,
    SingerName STRING(1024),
    SYNONYM (Artists)
) PRIMARY KEY (SingerId)