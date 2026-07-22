SELECT * FROM GRAPH_TABLE (FinGraph MATCH (n:Person) COLUMNS (n.*, 1 AS number))
