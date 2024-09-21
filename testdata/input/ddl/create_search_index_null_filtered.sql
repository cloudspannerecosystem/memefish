CREATE SEARCH INDEX AlbumsIndex
ON Albums(AlbumTitle_Tokens)
STORING(Genre)
WHERE Genre IS NOT NULL