create table types (
  b bool,
  i int64,
  f32 float32,
  f float64,
  d date,
  t timestamp,
  s string(256),
  smax string(max),
  bs bytes(256),
  bsmax bytes(max),
  ab array<bool>,
  abs array<bytes(max)>
) primary key (i)
