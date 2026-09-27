GRAPH FinGraph
MATCH (p:Person)
RETURN VALUE {
  MATCH (p)-[:Owns]->(a:Account)
  RETURN a.id AS id
  LIMIT 1
} AS account_id
