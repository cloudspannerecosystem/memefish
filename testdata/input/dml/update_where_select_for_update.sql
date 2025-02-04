update Albums set MarketingBudget = MarketingBudget + 100
where (SingerId, AlbumId) = (select as struct SingerId, AlbumId from Albums where AlbumTitle like "A%" for update)