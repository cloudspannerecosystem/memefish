ALTER TABLE metrics ALTER COLUMN doubled INT64 AS (source * 2) DEFAULT (0)
