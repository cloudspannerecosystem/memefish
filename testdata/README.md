# About testdata

If you place `.sql` files containing GoogleSQL in the correct location in this directory,
they will be automatically tested.

- inputs
  - ddl: input of `ParseDDL()` and `ParseStatement()`
  - dml: input of `ParseDML()` and `ParseStatement()`
  - expr: input of `ParseExpr()`
  - query: input of `ParseQuery()` and `ParseStatement()`
  - gql: input of `ParseGQLQuery()` and `ParseStatement()`
  - gql_graph_pattern: input of `ParseGQLGraphPattern()`
  - statement: input of `ParseStatement()`
- snapshots: the snapshot files of the parse results (AST and unparsed SQL) of `inputs`
- fuzz: the fuzzing corpus of `Fuzz*` functions in `fuzz_test.go`

## Snapshots

You can use this command in your project root to automatically update the snapshot files in `testdata/snapshots`.

```
$ make update-snapshots
```

Note: You should carefully check the diff when committing the snapshot files in `testdata/snapshots`.

## Tips

You can use ZetaSQL to check if it's a valid GoogleSQL query.

* This example requires to preload ZetaSQL docker container. See [Run with Docker](https://github.com/google/zetasql/tree/master?tab=readme-ov-file#run-with-docker).
* Currently, it is useful for query, DML, and expressions because DDL of Spanner GoogleSQL dialect is not compatible to ZetaSQL.

```sh
# statement
$ docker run --rm --platform linux/amd64 zetasql execute_query --product_mode=external --mode=parse,unparse "$(cat testdata/inputs/query/pipe_from_where_select_distinct.sql)"
# or expression
$ docker run --rm --platform linux/amd64 zetasql execute_query --product_mode=external --sql_mode=expression --mode=parse,unparse "$(cat testdata/inputs/expr/array_literal_empty_with_types.sql)"
```
