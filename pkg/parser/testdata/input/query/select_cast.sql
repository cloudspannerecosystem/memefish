select cast(1 as INT64), cast((struct(), 1, [2, 3], ["4", "5"]) as struct<struct<>, x int64, y array<int64>, z array<string>>)
from x tablesample BERNOULLI (cast(0.1 as float64) percent),
     y tablesample BERNOULLI (cast(1 as int64) rows),
     z tablesample BERNOULLI (cast(@param as int64) rows)
limit cast(1 as INT64) offset cast(@foo as INT64)
