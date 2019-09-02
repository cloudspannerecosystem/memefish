SELECT
  SingerID
FROM
  Singers
GROUP BY
  SingerID
HAVING
  SingerID = 1
