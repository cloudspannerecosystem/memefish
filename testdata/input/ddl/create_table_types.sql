-- https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#data_types
create table types (
  b bool,
  i int64,
  f32 float32,
  f float64,
  s string(256),
  sh string(0x100),
  smax string(max),
  bs bytes(256),
  bh bytes(0x100),
  bsmax bytes(max),
  j json,
  d date,
  t timestamp,
  ab array<bool>,
  abs array<bytes(max)>,
  af32vl array<float32>(vector_length=>128),
  p ProtoType,
  p_quoted `ProtoType`,
  p_path examples.ProtoType,
  p_partly_quoted_path examples.shipping.`Order`,
  p_fully_quoted_path `examples.shipping.Order`,
  ap ARRAY<ProtoType>,
  ap_quoted ARRAY<`ProtoType`>,
  ap_path ARRAY<examples.ProtoType>,
  ap_partly_quoted_path ARRAY<examples.shipping.`Order`>,
  ap_fully_quoted_path ARRAY<`examples.shipping.Order`>,
) primary key (i)
