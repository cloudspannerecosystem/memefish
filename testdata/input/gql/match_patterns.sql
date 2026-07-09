GRAPH FinGraph
-- Multiple path patterns in one MATCH
MATCH
  (src:Account)-[t1:Transfers]->(mid:Account)-[t2:Transfers]->(dst:Account),
  (mid)<-[:Owns]-(p:Person)
-- Path variable and property filter
MATCH p=(account:Account {is_blocked: false})-[transfer:Transfers]-(dst:Account)
-- Search prefixes
MATCH ANY SHORTEST (a:Account {is_blocked:true})-[t:Transfers]->{1, 4} (b:Account)
MATCH ANY CHEAPEST (a)-[e:Transfer COST e.amount]->{1,3}(b)
-- Subpath modes
MATCH WALK (a1:Account)-[t1:Transfers]->(a2:Account)
MATCH ACYCLIC (a1:Account)-[t1:Transfers]->(a2:Account)
MATCH TRAIL (a1:Account)-[t1:Transfers]->(a2:Account)
-- Subpaths with WHERE and hints
MATCH ((n) WHERE n.id = 1)
MATCH ( @{a=1} (n) )
MATCH (-[e:LocatedIn]->(p:Person)->(c:City) WHERE p.id = e.id)
MATCH (-[e:LocatedIn]->((p:Person)->(c:City)) WHERE p.id = e.id)
-- Consecutive subpaths
MATCH ((p:Person))(-[e:LocatedIn]->(c:City))(-[s:StudyAt]->(u:School))
-- Quantified subpath
MATCH (source_person:Person)((p:Person)-[k:Knows]->(f:Person)){, 4}(dest_person:Person)
RETURN 1
