GRAPH FinGraph
MATCH (a)-[e]->(b)
WHERE a IS SOURCE OF e
RETURN a.id
