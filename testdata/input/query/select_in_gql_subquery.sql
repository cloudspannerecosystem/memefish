SELECT 'Dana' IN {
  GRAPH FinGraph
  MATCH (p:Person)-[o:Owns]->(a:Account)
  RETURN p.name
} AS results