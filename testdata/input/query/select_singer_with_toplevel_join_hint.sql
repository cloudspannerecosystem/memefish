@{FORCE_JOIN_ORDER=TRUE}
SELECT
  *
FROM
  Singers A
  LEFT OUTER JOIN
  Singers B
  ON A.SingerID = B.SingerID
