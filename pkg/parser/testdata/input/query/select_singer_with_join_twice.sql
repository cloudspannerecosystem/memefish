SELECT
  *
FROM
  Singers A
  JOIN
  Singers B
  ON A.SingerID = B.SingerID
  INNER JOIN
  Singers C
  ON A.SingerID = C.SingerID
