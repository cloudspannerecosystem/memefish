-- https://cloud.google.com/spanner/docs/reference/standard-sql/operators#with_expression
SELECT WITH(a AS '123',       -- a is '123'
    b AS CONCAT(a, '456'),    -- b is '123456'
    c AS '789',               -- c is '789'
    CONCAT(b, c)) AS result   -- b + c is '123456789'