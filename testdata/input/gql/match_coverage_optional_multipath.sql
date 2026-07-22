GRAPH FinGraph
-- OPTIONAL MATCH
OPTIONAL MATCH (n:Person)
-- MATCH with Hints and multiple paths with top-level WHERE
MATCH @{join_method=hash_join} (src:Account)-[:Transfers]->(dst:Account), (src)<-[:Owns]-(p:Person)
WHERE p.id > 0
RETURN *
