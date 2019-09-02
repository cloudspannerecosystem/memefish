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
  OR SingerID BETWEEN 1 AND 3
  OR SingerID NOT BETWEEN 1 AND 3
  OR FirstName LIKE "%a"
  OR FirstName NOT LIKE "%a"
  OR NULL IS NULL
  OR NULL IS NOT NULL
  OR (SingerID = 1) IS TRUE
  OR (SingerID = 1) IS NOT TRUE
  OR (SingerID = 1) IS FALSE
  OR (SingerID = 1) IS NOT FALSE
