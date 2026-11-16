-- an ORDER BY column named interleave before the INTERLEAVE IN clause
CREATE SEARCH INDEX AlbumsIndex
ON Albums(AlbumTitle_Tokens)
PARTITION BY SingerId
ORDER BY interleave DESC, INTERLEAVE IN Singers