SELECT
  *
FROM
  Singers A
  HASH JOIN
  Singers B
  ON A.SingerID = B.SingerID