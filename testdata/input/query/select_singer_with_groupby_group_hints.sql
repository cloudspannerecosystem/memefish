SELECT
  FirstName, BirthDate
FROM
  Singers
GROUP @{GROUP_METHOD=HASH_GROUP} BY
  FirstName, BirthDate
