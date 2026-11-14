CREATE SEQUENCE MySequence BIT_REVERSED_POSITIVE START COUNTER WITH 50 OPTIONS (
    skip_range_min = 1,
    skip_range_max = 1000)