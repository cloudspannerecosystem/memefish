GRAPH FinGraph
MATCH (p:Person)
RETURN p.name, p.id
OFFSET 1
LIMIT 1