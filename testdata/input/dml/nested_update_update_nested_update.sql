UPDATE Singers s
SET (UPDATE s.AlbumInfo.Song so
     SET (UPDATE so.Chart c
          SET c.rank = 2
          WHERE c.chartname = "Galaxy Top 100")
     WHERE so.songtitle = "Bonus Track")
WHERE s.SingerId = 5