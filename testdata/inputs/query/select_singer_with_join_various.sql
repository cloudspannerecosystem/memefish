SELECT
  *
FROM
  Singers A
  JOIN
  Singers B
  ON A.SingerID = B.SingerID
  INNER JOIN
  Singers C
  ON A.SingerID = C.SingerID
  CROSS JOIN
  Singers D
  FULL JOIN
  Singers E
  ON A.SingerID = E.SingerID
  FULL OUTER JOIN
  Singers F
  ON A.SingerID = F.SingerID
  LEFT JOIN
  Singers G
  ON A.SingerID = G.SingerID
  LEFT OUTER JOIN
  Singers H
  ON A.SingerID = H.SingerID
  RIGHT JOIN
  Singers I
  ON A.SingerID = I.SingerID
  RIGHT OUTER JOIN
  Singers J
  ON A.SingerID = J.SingerID
