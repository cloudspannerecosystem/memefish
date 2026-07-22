SELECT name, number FROM GRAPH_TABLE (FinGraph MATCH (n:Person) COLUMNS (n.name AS name, 1 AS number))
