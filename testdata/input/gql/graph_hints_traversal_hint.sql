GRAPH FinGraph
MATCH
  (p:Person {id: 1})-[:Owns]->(a:Account),                   -- path pattern 1
  @{JOIN_METHOD=HASH_JOIN, HASH_JOIN_BUILD_SIDE=BUILD_RIGHT} -- traversal hint
  (a:Account)-[e:Transfers]->(c:Account)                     -- path pattern 2
RETURN c.id