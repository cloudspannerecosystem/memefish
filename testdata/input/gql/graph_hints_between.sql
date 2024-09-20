GRAPH FinGraph
MATCH
  (p:Person {id: 1})-[e:Owns]->
  @{JOIN_METHOD=APPLY_JOIN}
  ((a:Account)-[s:Transfers]->(oa:Account))
RETURN oa.id