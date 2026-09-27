SELECT
  *
FROM (
  SELECT
    *
  FROM
    Singers
  WHERE
    SingerID = 1
) as S
