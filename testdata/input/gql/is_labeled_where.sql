GRAPH FinGraph
MATCH (a)
WHERE a IS LABELED Account | Person
RETURN a.id
