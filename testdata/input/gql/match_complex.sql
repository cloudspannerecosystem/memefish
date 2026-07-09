GRAPH g
MATCH ANY SHORTEST (a:Person {name: "Alice"})-[:KNOWS]->{1, 5}(b:Person)
WHERE a.age > 20
RETURN a, b
