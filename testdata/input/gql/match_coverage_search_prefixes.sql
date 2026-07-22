GRAPH FinGraph
-- Search Prefixes
MATCH ALL (p1:Person)-[:Knows]->+(f1:Person)
MATCH ANY (p2:Person)-[:Knows]->+(f2:Person)
MATCH ANY SHORTEST (p3:Person)-[:Knows]->+(f3:Person)
MATCH ANY CHEAPEST (p4:Person)-[:Knows]->+(f4:Person)
RETURN *
