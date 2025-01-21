UPDATE Singers s
SET
  (DELETE FROM s.SingerInfo.Residence r WHERE r.City = 'Seattle'),
  (UPDATE s.Albums.Song song SET song.songtitle = 'No, This Is Rubbish' WHERE song.songtitle = 'This Is Pretty Good'),
  (INSERT s.Albums.Song VALUES ("songtitle: 'The Second Best Song'"))
WHERE SingerId = 3 AND s.Albums.title = 'Go! Go! Go!'