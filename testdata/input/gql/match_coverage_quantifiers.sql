GRAPH FinGraph
-- Quantifiers
MATCH (aq1)-[:Knows]->* (bq1)
MATCH (aq2)-[:Knows]->+ (bq2)
MATCH (aq4)-[:Knows]->{2} (bq4)
MATCH (aq5)-[:Knows]->{2,} (bq5)
MATCH (aq6)-[:Knows]->{,5} (bq6)
MATCH (aq7)-[:Knows]->{2,5} (bq7)
-- Subpath Patterns with everything
MATCH ( @{a=1} WALK (asub)-[esub]->(bsub) WHERE asub.id > 0 ){2,3}
RETURN *
