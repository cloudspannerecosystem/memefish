create index foo_bar on foo (
  bar,
  baz desc,
)
storing (qux)
where bar is not null and baz is not null
, interleave in parent
options (locality_group = 'spill_to_hdd')
