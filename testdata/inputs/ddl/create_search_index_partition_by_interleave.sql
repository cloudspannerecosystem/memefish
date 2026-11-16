CREATE SEARCH INDEX AlbumsIndex
ON Albums(AlbumTitle_Tokens)
PARTITION BY SingerId, AlbumId
, INTERLEAVE IN Singers