GRAPH FinGraph
MATCH (p:Person)
RETURN p.name IN {
  MATCH (x:Person)
  WHERE x.id < 10
  RETURN x.name
} AS flag
