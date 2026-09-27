CREATE SEQUENCE IF NOT EXISTS MySequence OPTIONS (
    sequence_kind='bit_reversed_positive',
    skip_range_min = 1,
    skip_range_max = 1000,
    start_with_counter = 50)
