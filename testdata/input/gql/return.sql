GRAPH FinGraph
MATCH (:Account)-[:Transfers]->(account:Account)
RETURN account, COUNT(*) AS num_incoming_transfers
GROUP BY account

NEXT

MATCH (account:Account)<-[:Owns]-(owner:Person)
RETURN owner.name AS owner_name, num_incoming_transfers
ORDER BY num_incoming_transfers DESC