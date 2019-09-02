select * from ((((select 1 A union all (select 2)) union distinct (select 1)) limit 1) JOIN (select 1 A, 2 B) USING (A))
