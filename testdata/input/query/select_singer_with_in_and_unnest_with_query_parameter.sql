SELECT
  *
FROM
  Singers
WHERE
  SingerId IN UNNEST(@singerIDs)
