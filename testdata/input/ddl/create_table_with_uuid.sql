create table foo (
  id uuid not null default (new_uuid()),
  name string(max)
) primary key (id)
