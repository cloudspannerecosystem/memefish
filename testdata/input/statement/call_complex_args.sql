-- https://github.com/google/zetasql/blob/a516c6b26d183efc4f56293256bba92e243b7a61/zetasql/parser/testdata/call.test#L92C1-L93C1
call myprocedure(TABLE my.table, (SELECT * FROM my.another_table), mytvf(1, 2))