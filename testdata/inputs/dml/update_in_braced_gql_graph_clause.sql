UPDATE Account
SET is_active = TRUE
WHERE id IN {
  GRAPH FinGraph
  MATCH (a)
  RETURN a.id
}
