CREATE INDEX SingersByFirstLastName ON Singers(FirstName, LastName)
  OPTIONS (locality_group = 'spill_to_hdd')