create or replace view singernames
sql security invoker
as select
    singers.singerid as singerid,
    singers.firstname || ' ' || singers.lastname as name
from singers
