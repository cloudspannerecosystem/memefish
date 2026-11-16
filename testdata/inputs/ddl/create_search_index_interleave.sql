-- https://cloud.google.com/spanner/docs/full-text-search/search-indexes#define_an_interleaved_search_index
CREATE SEARCH INDEX AlbumsIndex
ON Albums(AlbumTitle_Tokens)
PARTITION BY SingerId,
INTERLEAVE IN Singers
OPTIONS (sort_order_sharding = true)