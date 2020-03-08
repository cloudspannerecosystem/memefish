select 1 + 2, 1 - 2,
       1 * 2, 2 / 2,
       +1++1, -1+-1,
       +1.2, -3.4,
       ~1 ^ ~1,
       1 ^ 2, 2 & 1, 2 | 1,
       1 << 2, 2 >> 1,
       foo.bar * +foo.bar * -foo.bar,
       (select 1 `1`).1,
       NOT NOT true,
       [1, 2, 3][offset(1)],
       [1, 2, 3][`offset`(1)],
       [1, 2, 3][ordinal(1)],
       case
       when 1 = 1 then "1 = 1"
       else            "else"
       end,
       case 1
       when 1 then "1"
       when 2 then "2"
       else        "other"
       end,
       date_add(date "2019-09-01", interval 5 day),
       timestamp_add(timestamp "2019-09-01 08:11:22", interval 5 hour),
       1 in (1, 2, 3),
       2 in unnest([1, 2, 3]),
       3 in (select 1 union all select 2 union all select 3),
       [1] || [2]
