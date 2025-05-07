GRAPH FinGraph
MATCH (p:Person)-[:Owns]->(account:Account)
RETURN
 p.name, account.id AS account_id,
 ARRAY {
   MATCH (a:Account)-[transfer:Transfers]->(:Account)
   WHERE a = account
   RETURN transfer.amount AS transfers
 } AS transfers