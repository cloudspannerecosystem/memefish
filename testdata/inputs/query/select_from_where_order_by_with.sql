-- https://cloud.google.com/spanner/docs/full-text-search/ranked-search#score_multiple_columns
SELECT AlbumId
FROM Albums
WHERE SEARCH(Title_Tokens, @p1) AND SEARCH(Studio_Tokens, @p2)
ORDER BY WITH(
  TitleScore AS SCORE(Title_Tokens, @p1) * @titleweight,
  StudioScore AS SCORE(Studio_Tokens, @p2) * @studioweight,
  DaysOld AS (UNIX_MICROS(CURRENT_TIMESTAMP()) - ReleaseTimestamp) / 8.64e+10,
  FreshnessBoost AS (1 + @freshnessweight * GREATEST(0, 30 - DaysOld) / 30),
  PopularityBoost AS (1 + IF(HasGrammy, @grammyweight, 0)),
  (TitleScore + StudioScore) * FreshnessBoost * PopularityBoost)
LIMIT 2