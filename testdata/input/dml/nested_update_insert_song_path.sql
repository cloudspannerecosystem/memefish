UPDATE Singers s
SET (INSERT s.AlbumInfo.Song
     VALUES ('''songtitle: 'Bonus Track', length:180''')),
    s.Albums.tracks = 16
WHERE s.SingerId = 5 and s.AlbumInfo.title = "Fire is Hot"