select
    [1, 2, 3][offset(1)],
    [1, 2, 3][ordinal(1)],
    [1, 2, 3][safe_offset(1)],
    [1, 2, 3][ordinal(1)],
    [1, 2, 3][1],
    STRUCT(1, 2, 3)[offset(1)],
    STRUCT(1, 2, 3)[ordinal(1)],
    STRUCT(1, 2, 3)[safe_offset(1)],
    STRUCT(1, 2, 3)[ordinal(1)],
    STRUCT(1, 2, 3)[1],
    JSON '[1, 2, 3]'[1],
    JSON '{"a": 1, "b": 2, "c": 3}'['a']
