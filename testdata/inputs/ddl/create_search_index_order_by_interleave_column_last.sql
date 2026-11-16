-- a column named INTERLEAVE as the last ORDER BY key, no INTERLEAVE IN clause
CREATE SEARCH INDEX AlbumsIndex
ON Albums(AlbumTitle_Tokens)
PARTITION BY SingerId
ORDER BY ReleaseTimestamp DESC, INTERLEAVE