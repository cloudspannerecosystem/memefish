GRAPH FinGraph
MATCH (p:Person)
RETURN EXISTS { (p)-[:Owns]->(:Account) } AS has_account
