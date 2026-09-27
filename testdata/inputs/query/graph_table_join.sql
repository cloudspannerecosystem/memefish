SELECT a.id, g.v
FROM Accounts AS a
JOIN GRAPH_TABLE (FinGraph LET x = 1 RETURN x AS v) AS g
ON a.id = g.v
