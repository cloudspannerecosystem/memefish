SELECT * FROM (
  SELECT * FROM Singers
  |> SET FirstName = UPPER(FirstName), `Last Name` = "unknown",
)
