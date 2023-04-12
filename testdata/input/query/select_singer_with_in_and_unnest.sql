SELECT
  *
FROM
  Singers
WHERE
  SingerId IN UNNEST(ARRAY[1, 2, 3])
