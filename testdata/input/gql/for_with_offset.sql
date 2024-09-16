GRAPH FinGraph
MATCH (p:Person)-[o:Owns]->(a:Account)
FOR element in ["all","some"] WITH OFFSET
RETURN p.name, element as alert_type, offset
ORDER BY p.name, element, offset
