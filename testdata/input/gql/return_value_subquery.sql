GRAPH FinGraph
RETURN VALUE {
  MATCH (p:Person)
  WHERE p.name LIKE '%e%'
  RETURN p.name
  LIMIT 1
} AS results