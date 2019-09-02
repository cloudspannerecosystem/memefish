SELECT
  *
FROM
  Singers
WHERE
  SingerID = @singerID
  AND @singerID = SingerID
