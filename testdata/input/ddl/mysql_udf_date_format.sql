CREATE OR REPLACE FUNCTION mysql.DATE_FORMAT(d TIMESTAMP, format STRING) RETURNS STRING AS (
  SAFE.FORMAT_TIMESTAMP(
    CASE
      WHEN SAFE.REGEXP_CONTAINS(format, '(^|[^%])%[cDfhciMrsuVWXx]') THEN ERROR(
        SAFE.CONCAT(
          '%',
          SAFE.REGEXP_EXTRACT(format, '(^|[^%])%[cDfhciMrsuVWXx]'),
          ' format specifier is not supported'))
      ELSE format
    END,
    d)
)