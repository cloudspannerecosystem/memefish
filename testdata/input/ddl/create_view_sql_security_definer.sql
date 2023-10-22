create view singernames
sql security definer
as select
    singers.singerid as singerid,
    singers.firstname || ' ' || singers.lastname as name
from singers
