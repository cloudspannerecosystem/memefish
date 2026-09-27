CREATE TABLE Singers (
  SingerId   INT64 NOT NULL,
  FirstName  STRING(1024),
  LastName   STRING(1024),
  Awards     ARRAY<STRING(MAX)> OPTIONS (locality_group = 'spill_to_hdd')
) PRIMARY KEY (SingerId), OPTIONS (locality_group = 'ssd_only')