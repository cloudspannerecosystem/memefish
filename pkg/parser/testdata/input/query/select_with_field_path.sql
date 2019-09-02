SELECT
  A.x,
  A.y,
  A.z.a,
  A.z.b
FROM
  UNNEST(
    ARRAY(
      SELECT AS STRUCT
        x,
        y,
        z
      FROM
        UNNEST(ARRAY<STRUCT<x INT64, y STRING, z STRUCT<a INT64, b INT64>>>[(1, 'foo', (2, 3)), (3, 'bar', (4, 5))])
    )
  ) AS A
WHERE A.z.a = 2
