GRAPH FinGraph
MATCH (p:Person)
FOR element in [] WITH OFFSET AS off
RETURN p.name, element, off