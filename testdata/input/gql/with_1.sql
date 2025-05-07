GRAPH FinGraph
MATCH (:Account)-[:Transfers]->(account:Account)
WITH account, COUNT(*) AS num_incoming_transfers GROUP BY account
RETURN account.id AS account_id, num_incoming_transfers