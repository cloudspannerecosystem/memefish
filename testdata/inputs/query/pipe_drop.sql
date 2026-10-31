SELECT * FROM (
  SELECT * FROM Singers
  |> DROP FirstName, `Last Name`,
)
