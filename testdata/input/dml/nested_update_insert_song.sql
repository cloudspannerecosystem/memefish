UPDATE Singers s
SET (INSERT s.AlbumInfo.Song(Song)
     VALUES ("songtitle: 'Bonus Track', length: 180"))
WHERE s.SingerId = 5 AND s.AlbumInfo.title = "Fire is Hot"