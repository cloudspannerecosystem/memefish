GRAPH FinGraph
MATCH (n:Person)
OPTIONAL MATCH (n:Person)-[:Owns]->(a:Account {is_blocked: TRUE})
RETURN n.name, a.id AS blocked_account_id