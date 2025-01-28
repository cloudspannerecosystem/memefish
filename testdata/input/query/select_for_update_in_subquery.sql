SELECT s.SingerId, s.FirstName, s.LastName FROM Singers AS s
JOIN
(SELECT SingerId FROM Albums WHERE MarketingBudget > 100000 FOR UPDATE) AS a
ON a.SingerId = s.SingerId
