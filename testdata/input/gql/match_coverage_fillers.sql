GRAPH FinGraph
-- Element Filler with Variable, Properties, WHERE, and COST
MATCH path = (ne:Person {age: 20})-[ee:Knows WHERE ee.since > 2020 COST ee.since]->(me)
RETURN *
NEXT
MATCH (ne:Person WHERE ne.name = 'Alice')
RETURN *
