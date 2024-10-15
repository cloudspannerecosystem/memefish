-- Should
CREATE SEARCH INDEX AlbumsIndexFull
ON Albums(Title_Tokens, Studio_Tokens)
STORING(Genre)
PARTITION BY SingerId
ORDER BY ReleaseTimestamp DESC
WHERE Genre IS NOT NULL AND ReleaseTimestamp IS NOT NULL
, INTERLEAVE IN Singers
OPTIONS(sort_order_sharding=true)