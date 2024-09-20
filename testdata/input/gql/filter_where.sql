GRAPH FinGraph
MATCH (p:Person)-[o:Owns]->(a:Account)
FILTER WHERE p.birthday < '1990-01-10'
RETURN p.name