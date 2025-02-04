@{pdml_max_parallelism=1}
insert into foo@{force_index=_base_table} (foo, bar, baz)
values (1, 2, 3),
       (4, 5, 6)