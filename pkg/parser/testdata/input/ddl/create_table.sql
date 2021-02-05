create table foo (
  foo int64,
  bar float64 not null,
  baz string(255) not null options(allow_commit_timestamp = null),
  qux string(255) not null as (concat(baz, "a")) stored,
  foreign key (foo) references t2 (t2key1),
  constraint fkname foreign key (foo, bar) references t2 (t2key1, t2key2)
) primary key (foo, bar)
