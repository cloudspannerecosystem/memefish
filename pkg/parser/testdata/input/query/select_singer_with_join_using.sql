SELECT
  *
FROM
  Singers A
  LEFT OUTER JOIN
  Singers B
  USING (SingerID, FirstName)
