GRAPH FinGraph
-- Path Modes and Subpaths
MATCH TRAIL (a)->{2,}(b)
RETURN *
NEXT
MATCH TRAIL ((a)-[e]->(b)){0, }
RETURN *
NEXT
MATCH ANY SHORTEST (TRAIL ->{1,4})
RETURN *
