create table foo (
  foo int64,
  bar float64 not null,
  baz string(255) not null options(allow_commit_timestamp = null)
) primary key (foo, bar)
