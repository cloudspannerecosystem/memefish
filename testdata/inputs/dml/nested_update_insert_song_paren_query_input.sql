-- In nested query, column list is optional so there is ambiguity between parenthesized query input and column list.
-- I believe Spanner hasn't yet supported this kind of query, but it can be parsed.
UPDATE Singers s
SET (INSERT s.AlbumInfo.Song (SELECT AS VALUE CAST("songtitle: 'The Second Best Song'" AS googlesql.example.Album.Song)))
WHERE TRUE