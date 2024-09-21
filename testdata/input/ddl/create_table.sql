create table foo (
  foo int64,
  bar float64 not null,
  baz string(255) not null options(allow_commit_timestamp = null),
  qux string(255) not null as (concat(baz, "a")) stored,
  foreign key (foo) references t2 (t2key1),
  foreign key (bar) references t2 (t2key2) on delete cascade,
  foreign key (baz) references t2 (t2key3) on delete no action,
  constraint fkname foreign key (foo, bar) references t2 (t2key1, t2key2),
  check (foo > 0),
  constraint cname check (bar > 0),
  quux json,
  corge timestamp not null default (current_timestamp()),
) primary key (foo, bar)
