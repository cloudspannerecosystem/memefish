GRAPH FinGraph
-- Label Expressions
MATCH (nl1:Person|Employee)
MATCH (nl2:Person&Employee)
MATCH (nl3:!Person)
MATCH (nl4:%)
MATCH (nl5:(Person|Employee)&!Student)
MATCH (nl6 IS Person)
MATCH (nl7:Person)
RETURN *
