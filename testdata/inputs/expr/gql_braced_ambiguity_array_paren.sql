-- Regression: ARRAY(...) SQL subquery must still parse after ARRAY { gql } support.
ARRAY(SELECT 1)