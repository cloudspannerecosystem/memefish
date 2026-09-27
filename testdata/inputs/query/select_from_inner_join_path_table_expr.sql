-- https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#correlated_join
SELECT A.name, item
FROM
  UNNEST(
    [
      STRUCT(
        'first' AS name,
        [1, 2, 3, 4] AS items),
      STRUCT(
          'second' AS name,
        [] AS items)]) AS A
    INNER JOIN
  A.items AS item
