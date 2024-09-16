GRAPH FinGraph
MATCH ANY SHORTEST (a:Account)-[t:Transfers]->{1, 4} (b:Account)
WHERE a.is_blocked
LET total = SUM(t.amount)
RETURN a.id AS a_id, total, b.id AS b_id