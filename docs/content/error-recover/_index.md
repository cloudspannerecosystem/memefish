---
date: 2024-12-20 00:00:00 +0900
title: "Error recovering"
weight: 2
---

Since v0.1.0, `memefish.ParseXXX` methods returns AST node(s) even if an error is reproted.
That is, if we try to parse incomplete SQL such as:

```sql
SELECT (1 +) + (* 2)
```

Then, the following two errors are reported:

```sql
syntax error: :1:12: unexpected token: )
  1|  SELECT (1 +) + (* 2)
   |             ^
syntax error: :1:17: unexpected token: *
  1|  SELECT (1 +) + (* 2)
   |                  ^
```

Hoever, the AST is also returned:

```go {hl_lines=["10-31","36-57"]}
&ast.QueryStatement{
  Query: &ast.Select{
    Results: []ast.SelectItem{
      &ast.ExprSelectItem{
        Expr: &ast.BinaryExpr{
          Op:   "+",
          Left: &ast.ParenExpr{
            Lparen: 7,
            Rparen: 11,
            Expr:   &ast.BadExpr{
              BadNode: &ast.BadNode{
                NodePos: 8,
                NodeEnd: 11,
                Tokens:  []*token.Token{
                  &token.Token{
                    Kind: "<int>",
                    Raw:  "1",
                    Base: 10,
                    Pos:  8,
                    End:  9,
                  },
                  &token.Token{
                    Kind:  "+",
                    Space: " ",
                    Raw:   "+",
                    Pos:   10,
                    End:   11,
                  },
                },
              },
            },
          },
          Right: &ast.ParenExpr{
            Lparen: 15,
            Rparen: 19,
            Expr:   &ast.BadExpr{
              BadNode: &ast.BadNode{
                NodePos: 16,
                NodeEnd: 19,
                Tokens:  []*token.Token{
                  &token.Token{
                    Kind: "*",
                    Raw:  "*",
                    Pos:  16,
                    End:  17,
                  },
                  &token.Token{
                    Kind:  "<int>",
                    Space: " ",
                    Raw:   "2",
                    Base:  10,
                    Pos:   18,
                    End:   19,
                  },
                },
              },
            },
          },
        },
      },
    },
  },
}
```

Thus, the places where the error occurred are filled with the `ast.BadXXX` nodes (`ast.BadExpr` in this example).

## How méméfish performs error recovery

This section explains how méméfish performs error recovery.

In méméfish, a *recovery point* is set when parsing a syntax where some multiple types of AST nodes are the result.
For example, when parsing an parenthesized expression, the recovery point is set after the open parenthesis `(`.
If an error occurs in the parenthesized expression, the parser backtracks to the recovery point and skips the tokens until the parenthesized expression ends.
The skipped tokens are then collectively `ast.BadNode` and this node is wrapped up a specific `ast.BadXXX` node (e.g., `ast.BadExpr`).

```sql
SELECT (1 + 2 *)
               ^--- error point
        ^---------- recovery point
        |~~~~~| --- skipped tokens
```

Recovery points are set where:

- the beginning of statements, queries, DDLs, DMLs,
- the beginning of expressions (e.g., after an open parenthesis `(`, `SELECT`, `WHERE` etc.), and
- the beginning of types.

Token skipping is performed as follows.

- For `ast.Statement`, `ast.DDL`, and `ast.DML`,
  * skip tokens until a semicolon `;` appears.
- For `ast.QueryExpr`,
  * skip tokens until a semicolon `;` appears, or
  * skip tokens with counting the nest of parentheses `(`
    + until the closing symbol (`)`) appears at no nestings, or
    + until the symbol that is supposed to be the end of the expression (`UNION`, `INTERSECT`, `EXCEPT`) appears at no nestings.
- For `ast.Expr`,
  * skip tokens until a semicolon `;` appears, or
  * skip tokens with counting the nest of parentheses `(`, brackets `[`, `CASE` and `WHEN`
    + until the closing symbol (`)`, `]`, `END`, `THEN`) appears at no nestings or
    + until the symbol that is supposed to be the end of the expression (`,`, `AS`, `FROM`, `GROUP`, `HAVING`, `ORDER`, `LIMIT`, `OFFSET`, `AT`, `UNION`, `INTERSECT`, `EXCEPT`) appears at no nestings.
- For `ast.Type`,
  * skip tokens until the semicolon `;` or the closing parenthesis `)` appears, or
  * skip tokens with counting the nest of triangle brackets `<`
  * until the closing symbol (`>`) appears at no nestings.

Note that this skipping rules are just heuristics and may not be perfect.
In some cases, there is a possibility of skipping too many tokens.
