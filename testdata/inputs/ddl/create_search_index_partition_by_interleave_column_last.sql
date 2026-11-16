-- a column named interleave as the last PARTITION BY key, no INTERLEAVE IN clause
CREATE SEARCH INDEX AlbumsIndex
ON Albums(AlbumTitle_Tokens)
PARTITION BY SingerId, interleave
OPTIONS (sort_order_sharding = true)