GRAPH FinGraph
RETURN
  EXISTS @{JOIN_METHOD=HASH_JOIN} { ANY SHORTEST (a)-[]->(b) } AS any_shortest,
  EXISTS { TRAIL PATH (a)-[]->(b) } AS trail_path,
  VALUE @{a=1} { RETURN 1 AS value } AS hinted_value
