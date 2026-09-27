CREATE TABLE sch1.Singers (
    SingerId INT64 NOT NULL,
    FirstName STRING(1024),
    LastName STRING(1024),
    SingerInfo BYTES(MAX),
) PRIMARY KEY(SingerId)