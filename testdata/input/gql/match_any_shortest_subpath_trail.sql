GRAPH FinGraph
MATCH ANY SHORTEST (TRAIL ->{1,4})
RETURN COUNT(1) as num_paths