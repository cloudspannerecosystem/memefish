---
title: README
---

<p align="center">
  <img src="images/memefish.png" width="220px">
</p>

# méméfish

> méméfish is the foundation to analyze [Spanner][] [SQL][Spanner SQL].

[Spanner]: https://cloud.google.com/spanner/
[Spanner SQL]: https://cloud.google.com/spanner/docs/query-syntax

[![GoDoc Reference][godoc-badge]](https://godoc.org/github.com/cloudspannerecosystem/memefish/pkg)

## News

<table>
  <tr><th>ℹ️</th><td>

Since 2023/4/1, this repository has been moved from [MakeNowJust](https://github.com/makenowjust) to [cloudspannerecosystem](https://github.com/cloudspannerecosystem). You may need to migrate import paths from `github.com/MakeNowJust/memefish` to `github.com/cloudspannerecosystem/memefish` like:

```diff
 import (
-	"github.com/MakeNowJust/memefish/pkg/parser"
+	"github.com/cloudspannerecosystem/memefish/pkg/parser"
 )
```

  </td></tr>

  <tr><th>ℹ️</th><td>

Since 2023/4/8, the layout of this repository has been changed. Now, the old `parser` package has been moved to the top of the repository as the new `memefish` package, and sub-packages in the `pkg` directory are placed under the top.
  </td></tr>
</table>


## Features

- Parse Spanner SQL to AST
- Generate Spanner SQL from AST (unparse)

## Notice

This project is originally developed under "Expert team Go Engineer (Backend)" of [Mercari Summer Internship for Engineer 2019](https://mercan.mercari.com/articles/13497/).

## License

This project is licensed under MIT license.

[godoc-badge]: https://img.shields.io/badge/godoc-reference-black.svg?style=for-the-badge&colorA=%235272B4&logo=go&logoColor=white
