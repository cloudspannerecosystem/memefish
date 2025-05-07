GRAPH FinGraph
MATCH (n:Person|Account)
RETURN n.id, n.name, n.nick_name