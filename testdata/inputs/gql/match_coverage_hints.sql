GRAPH FinGraph
-- Hints coverage
MATCH (a:Person {name: "Alex"}), @{JOIN_METHOD=HASH_JOIN} (a)-[:Owns]->(acc:Account)
RETURN *
NEXT
MATCH (p:Person {id:1})-[e:Owns]->@{JOIN_METHOD=APPLY_JOIN}(a:Account)
RETURN *
NEXT
MATCH (p:Person {id: 1})@{JOIN_METHOD=APPLY_JOIN}-[e:Owns]->(a:Account)
RETURN *
NEXT
MATCH (a:Account {id:7})-[@{INDEX_STRATEGY=FORCE_INDEX_UNION} :Transfers]-(oa:Account)
RETURN *
NEXT
RETURN p.name, COUNT(*) AS num_account
GROUP @{GROUP_METHOD=HASH_GROUP} BY p.name
