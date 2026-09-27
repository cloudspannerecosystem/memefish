SELECT SAFE.SUBSTR('foo', 0, -2) AS safe_output UNION ALL
SELECT SAFE.SUBSTR('bar', 0, 2) AS safe_output