-- a key column named Interleave before the INTERLEAVE IN clause
CREATE SEARCH INDEX AlbumsIndex
ON Albums(AlbumTitle_Tokens)
PARTITION BY SingerId, Interleave, INTERLEAVE IN Singers