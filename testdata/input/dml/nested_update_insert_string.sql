UPDATE Singers s
SET (INSERT s.AlbumInfo.comments
     VALUES ("Groovy!"))
WHERE s.SingerId = 5 AND s.AlbumInfo.title = "Fire is Hot"