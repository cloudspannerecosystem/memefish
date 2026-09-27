SELECT a.AlbumId, a.Description
FROM Albums a
WHERE a.SingerId = 1 AND SEARCH(a.DescriptionTokens, 'classic albums', enhance_query => TRUE)