GRAPH FinGraph
MATCH (src:Account {id: 7})-[e:Transfers]->{1, 3}(dst:Account)
RETURN DISTINCT ARRAY_LENGTH(e) AS hops, dst.id AS destination_account_id