GRAPH FinGraph
MATCH CHEAPEST (a:Person)-[:Knows]->(b:Person)
RETURN a, b
