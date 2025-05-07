GRAPH FinGraph
RETURN 'Dana' IN {
  MATCH (p:Person)-[o:Owns]->(a:Account)
  RETURN p.name
} AS results