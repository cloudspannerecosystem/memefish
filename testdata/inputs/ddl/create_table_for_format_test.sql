create table if not exists foo (
    foo int64,
    bar float64 not null,
    baz string(255) not null options(allow_commit_timestamp = null),
    qux string(255) not null as (concat(baz, "a")) stored,
    foreign key (foo) references t2 (t2key1),
    constraint fkname foreign key (foo, bar) references t2 (t2key1, t2key2),
    check (foo > 0),
    constraint cname check (bar > 0),
    corge timestamp not null default (current_timestamp())
) primary key (foo),
  interleave in parent foobar,
  row deletion policy ( older_than ( baz, INTERVAL 30 DAY ) )
