GRAPH FinGraph
MATCH TRAIL (WALK (a1:Account)-[t1:Transfers]->{4}(a5:Account))
RETURN COUNT(1) as num_paths