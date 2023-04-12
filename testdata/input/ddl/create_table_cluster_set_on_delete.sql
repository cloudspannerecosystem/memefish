create table foo (
  foo int64
) primary key (foo),
  interleave in parent foobar
             on delete cascade