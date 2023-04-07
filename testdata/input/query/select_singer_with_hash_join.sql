SELECT
  *
FROM
  Singers A
  HASH JOIN
  Singers B
  ON A.SingerID = B.SingerID
  APPLY JOIN
  Singer C
  ON B.SingerID = C.SingerID
  LOOP JOIN
  Singer D
  ON C.SingerID = D.SingerID
