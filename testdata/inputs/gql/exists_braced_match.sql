GRAPH FinGraph
MATCH (p:Person)
RETURN EXISTS {
  MATCH (p)-[:Owns]->(a:Account)
} AS has_account
