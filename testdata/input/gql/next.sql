GRAPH FinGraph
MATCH (:Account)-[:Transfers]->(account:Account)
RETURN account, COUNT(*) AS num_incoming_transfers
GROUP BY account

NEXT

MATCH (account:Account)<-[:Owns]-(owner:Person)
RETURN
  account.id AS account_id, owner.name AS owner_name,
  num_incoming_transfers

NEXT

FILTER num_incoming_transfers < 2
RETURN account_id, owner_name
UNION ALL
RETURN "Bob" AS owner_name, 100 AS account_id