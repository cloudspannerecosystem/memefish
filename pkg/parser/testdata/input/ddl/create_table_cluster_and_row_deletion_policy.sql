create table foo (
  foo int64,
  bar int64,
  baz timestamp,
) primary key (),
  interleave in parent foobar,
  row deletion policy ( older_than ( baz, INTERVAL 30 DAY ) )
