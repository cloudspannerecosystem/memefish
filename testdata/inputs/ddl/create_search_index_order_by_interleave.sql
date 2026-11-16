CREATE SEARCH INDEX AlbumsIndex
ON Albums(AlbumTitle_Tokens)
PARTITION BY SingerId
ORDER BY ReleaseTimestamp DESC
, INTERLEAVE IN Singers
OPTIONS(sort_order_sharding=true)