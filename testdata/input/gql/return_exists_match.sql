GRAPH FinGraph
RETURN EXISTS {
  MATCH (p:Person)-[o:Owns]->(a:Account)
  WHERE p.Name LIKE 'D%'
} AS results