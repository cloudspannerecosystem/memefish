GRAPH FinGraph
MATCH
  (TRAIL (a1:Account)-[t1:Transfers]->{3}(a4:Account))
  -[t4:Transfers]->(a5:Account)
RETURN COUNT(1) as num_paths