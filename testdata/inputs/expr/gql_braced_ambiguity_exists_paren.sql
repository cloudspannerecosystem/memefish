-- Regression: EXISTS(...) SQL subquery must still parse after EXISTS { gql } support.
EXISTS @{JOIN_METHOD=HASH_JOIN} (SELECT 1)
