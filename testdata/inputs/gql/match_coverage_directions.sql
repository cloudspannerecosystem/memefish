GRAPH FinGraph
-- Edge Directions and Fillers
MATCH (ad1)-[:Knows]-(bd1)
MATCH (ad2)<-[:Knows]-(bd2)
MATCH (ad3)-[:Knows]->(bd3)
MATCH (ad5)-[]-(bd5)
MATCH (ad6)-(bd6)
MATCH (ad7)<-(bd7)
MATCH (ad8)->(bd8)
RETURN *
