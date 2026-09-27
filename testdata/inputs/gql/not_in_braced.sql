GRAPH FinGraph
MATCH (p:Person)
RETURN p.name NOT IN {
  MATCH (x:Person)
  RETURN x.name
} AS flag
