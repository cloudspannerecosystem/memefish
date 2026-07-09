GRAPH FinGraph
MATCH SHORTEST (a:Person)-[:Knows]->(b:Person)
RETURN a, b
