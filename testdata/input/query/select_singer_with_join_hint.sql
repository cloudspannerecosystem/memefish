SELECT
  *
FROM
  Singers A
  LEFT OUTER JOIN@{FORCE_JOIN_ORDER=TRUE}
  Singers B
  ON A.SingerID = B.SingerID
  JOIN@{JOIN_TYPE=HASH_JOIN}
  Singers C
  ON A.SingerID = C.SingerID
  JOIN@{JOIN_TYPE=APPLY_JOIN}
  Singers D
  ON A.SingerID = D.SingerID
  JOIN@{JOIN_TYPE=LOOP_JOIN}
  Singers E
  ON A.SingerID = E.SingerID
