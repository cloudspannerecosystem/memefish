create table foo (
  foo int64,
  bar int64,
  baz timestamp,
) primary key (),
  row deletion policy ( older_than ( baz, INTERVAL 30 DAY ) )
