SELECT
  *
FROM
  ComplexTable,
  UNNEST(ComplexTable.IntArray)
