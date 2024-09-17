GRAPH FinGraph
RETURN EXISTS {
  MATCH (p:Person)-[o:Owns]->(a:Account)
  WHERE p.Name LIKE 'D%'
  RETURN p.Name
  LIMIT 1
} AS results