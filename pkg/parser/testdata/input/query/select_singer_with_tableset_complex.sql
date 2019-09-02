SELECT * FROM Singers
UNION ALL
(
  SELECT * FROM Singers
  UNION DISTINCT
  (
    SELECT * FROM Singers
    INTERSECT ALL
    (
      SELECT * FROM Singers
      INTERSECT DISTINCT
      (
        SELECT * FROM Singers
        EXCEPT ALL
        (
          SELECT * FROM Singers
          EXCEPT DISTINCT
          SELECT * FROM Singers
        )
      )
    )
  )
)
