# About testdata

If you place `.sql` files containing GoogleSQL in the correct location in this directory,
they will be automatically tested.

- input
  - ddl: input of `ParseDDL()`
  - dml: input of `ParseDML()`
  - expr: input of `ParseExpr()`
  - query: input of `ParseQueryStatement()`
  - statement: input of `ParseStatement()`

You can use this command in your project root to automatically update `testdata/result`.

```
$ go test --update
```

Note: You should carefully check the diff when committing the contents of `testdata/result`.

## Tips

You can use ZetaSQL to check if it's a valid GoogleSQL query.

* This example requires to preload ZetaSQL docker container. See [Run with Docker](https://github.com/google/zetasql/tree/master?tab=readme-ov-file#run-with-docker).
* Currently, it is useful for query, DML, and expressions because DDL of Spanner GoogleSQL dialect is not compatible to ZetaSQL.

```sh
# statement
$ docker run --rm --platform linux/amd64 zetasql execute_query --product_mode=external --mode=parse,unparse "$(cat testdata/input/query/pipe_from_where_select_distinct.sql)"
# or expression
$ docker run --rm --platform linux/amd64 zetasql execute_query --product_mode=external --sql_mode=expression --mode=parse,unparse "$(cat testdata/input/expr/array_literal_empty_with_types.sql)" ```
```