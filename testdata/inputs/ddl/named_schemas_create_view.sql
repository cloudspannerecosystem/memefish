CREATE VIEW sch1.SingerView SQL SECURITY INVOKER
AS Select s.FirstName, s.LastName, s.SingerInfo
   FROM sch1.Singers AS s WHERE s.SingerId = 123456