SELECT
  *
FROM
  Singers
WHERE
  SingerID = 1
  OR SingerID < 1
  OR SingerID > 1
  OR SingerID <= 1
  OR SingerID >= 1
  OR SingerID != 1
  OR SingerID IN (1, 2, 3)
  OR SingerID NOT IN (1, 2, 3)
  OR FirstName LIKE "%a"
  OR FirstName NOT LIKE "%a"
