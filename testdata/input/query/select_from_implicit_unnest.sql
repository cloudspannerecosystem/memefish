SELECT a, off
FROM UNNEST([STRUCT<arr ARRAY<STRING>>(["foo"])]) AS t,
     t.arr AS a WITH OFFSET AS off