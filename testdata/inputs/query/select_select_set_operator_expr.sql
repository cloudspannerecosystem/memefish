select ((select 1) union all (select 2)) + 3,
       ((select 1) intersect all (select 1)) + 3,
       ((select 1) except all (select 1)) + 3
