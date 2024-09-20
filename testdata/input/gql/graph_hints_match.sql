GRAPH FinGraph
MATCH (p:Person {id: 1})-[:Owns]->(a:Account)
MATCH @{JOIN_METHOD=APPLY_JOIN}(a:Account)-[e:Transfers]->(oa:Account)
RETURN oa.id