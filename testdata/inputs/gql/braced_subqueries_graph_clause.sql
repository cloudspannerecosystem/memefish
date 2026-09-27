GRAPH FinGraph
MATCH (n)
RETURN
  EXISTS { GRAPH FinGraph MATCH (m) RETURN m } AS has_match,
  EXISTS { MATCH (m) RETURN m UNION ALL RETURN m } AS compound_match,
  EXISTS { MATCH (m) RETURN m NEXT RETURN m } AS chained_match,
  ARRAY { GRAPH FinGraph MATCH (m) RETURN m } AS matches,
  VALUE { GRAPH FinGraph MATCH (m) RETURN m } AS value,
  n IN { GRAPH FinGraph MATCH (m) RETURN m } AS contained
