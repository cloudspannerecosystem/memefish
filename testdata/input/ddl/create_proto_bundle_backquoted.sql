-- If you're using a protocol buffer type and any part of the type name is a Spanner reserved keyword,
-- enclose the entire protocol buffer type name in backticks.
CREATE PROTO BUNDLE (
       `examples.shipping.Order`,
       `examples.shipping.Order.Address`,
       `examples.shipping.Order.Item`)