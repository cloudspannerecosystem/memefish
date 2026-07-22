GRAPH FinGraph
MATCH p = ANY SHORTEST (a)->(b), shortest = (c), q = TRAIL PATHS (d), edge = (->)
RETURN *
