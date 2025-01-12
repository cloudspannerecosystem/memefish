-- https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#correlated_join
SELECT *
FROM
  Roster
    JOIN
  UNNEST(
      ARRAY(
        SELECT AS STRUCT *
      FROM PlayerStats
      WHERE PlayerStats.OpponentID = Roster.SchoolID
    )) AS PlayerMatches
  ON PlayerMatches.LastName = 'Buchanan'
