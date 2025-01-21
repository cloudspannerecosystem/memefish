UPDATE Singers s
SET (UPDATE s.AlbumInfo.Song so
    SET (INSERT INTO so.Chart
         VALUES ("chartname: 'Galaxy Top 100', rank: 5"))
    WHERE so.songtitle = "Bonus Track")
WHERE s.SingerId = 5