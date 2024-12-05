-- https://cloud.google.com/spanner/docs/reference/standard-sql/functions-reference#function_hints
SELECT
    SUBSTRING(CAST(x AS STRING), 2, 5) AS w,
    SUBSTRING(CAST(x AS STRING), 3, 7) AS y
FROM (SELECT SHA512(z) @{DISABLE_INLINE = TRUE} AS x FROM t)