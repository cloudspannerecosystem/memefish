SELECT p.x, q.x
FROM ML.PREDICT(MODEL MyModel, (SELECT 1 AS x)) @{remote_udf_max_rows_per_rpc=1} AS p
CROSS JOIN ML.PREDICT(MODEL MyModel, (SELECT 2 AS x)) q
