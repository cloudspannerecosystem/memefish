GRAPH FinGraph
MATCH ALL SHORTEST (a:Person)-[:Knows]->(b:Person)
MATCH ALL CHEAPEST (a2:Person)-[:Knows COST 1]->(b2:Person)
MATCH ANY 2 (a3:Person)-[:Knows]->(b3:Person)
MATCH SHORTEST 3 (a4:Person)-[:Knows]->(b4:Person)
MATCH CHEAPEST 1 (a5:Person)-[:Knows COST 1]->(b5:Person)
RETURN a, b
