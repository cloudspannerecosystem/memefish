SELECT
  ARRAY(
    (
      SELECT AS STRUCT
        *
      FROM Singers LIMIT 100
    )
  )
