SELECT
  *
FROM
  Singers A
  LEFT OUTER JOIN
  Singers B
  ON A.SingerID = B.SingerID
