GRAPH FinGraph
MATCH (source:Account)-[e:Transfers]->(destination:Account)
ORDER BY source.nick_name
LIMIT 3
RETURN source.nick_name