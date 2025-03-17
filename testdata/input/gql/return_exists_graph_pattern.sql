GRAPH FinGraph
RETURN EXISTS {
  (p:Person)-[o:Owns]->(a:Account) WHERE p.Name LIKE 'D%'
} AS results