CREATE QUEUE Q (
  id INT64 AS (1) STORED HIDDEN,
  ts TIMESTAMP OPTIONS (allow_commit_timestamp = true),
  Payload STRING(MAX) NOT NULL HIDDEN
) PRIMARY KEY (id, ts)
