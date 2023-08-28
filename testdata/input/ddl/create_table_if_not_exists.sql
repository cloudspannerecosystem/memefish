create table if not exists foo (
  foo int64,
  bar float64 not null,
) primary key (foo, bar)
