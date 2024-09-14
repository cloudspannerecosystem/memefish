GRAPH FinGraph
MATCH (source:Account)-[e:Transfers]->(destination:Account)
LET a = source
RETURN a.id AS a_id