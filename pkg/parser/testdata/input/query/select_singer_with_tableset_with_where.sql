SELECT * FROM Singers
UNION ALL
SELECT * FROM Singers
WHERE
  SingerId = 1
ORDER BY
  FirstName
