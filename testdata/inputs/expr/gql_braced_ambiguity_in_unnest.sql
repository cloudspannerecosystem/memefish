-- Regression: IN UNNEST(...) must still parse after IN { gql } support.
1 IN UNNEST([1, 2, 3])