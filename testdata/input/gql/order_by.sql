GRAPH FinGraph
MATCH (src_account:Account)-[transfer:Transfers]->(dst_account:Account)
ORDER BY transfer.amount DESC
LIMIT 3
RETURN src_account.id AS account_id, transfer.amount AS transfer_amount