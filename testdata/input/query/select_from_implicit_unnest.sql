SELECT *
FROM UNNEST([STRUCT<arr ARRAY<STRING>>(["foo"])]) AS t, t.arr