GRAPH FinGraph
MATCH (p:Person)
RETURN EXISTS {
  MATCH (p)-[:Owns]->(a:Account)
  RETURN a.id
} AS has_account
