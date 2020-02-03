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

[![GoDoc Reference][godoc-badge]](https://godoc.org/github.com/MakeNowJust/memefish/pkg)
[![Coverage][codecov-badge]](https://codecov.io/gh/MakeNowJust/memefish/branch/master)

## Features

- Parse Spanner SQL to AST
- Generate Spanner SQL from AST (unparse)
- Check expression type and semantics in SQL statement

## Notice

This project is originally developed under "Expert team Go Engineer (Backend)" of [Mercari Summer Internship for Engineer 2019](https://mercan.mercari.com/articles/13497/).

## License

This project is licensed under MIT license.

[godoc-badge]: https://img.shields.io/badge/godoc-reference-black.svg?style=for-the-badge&colorA=%235272B4&logo=go&logoColor=white
[codecov-badge]: https://img.shields.io/codecov/c/github/MakeNowJust/memefish/master.svg?style=for-the-badge&colorA=FF005E&logo=codecov&logoColor=white
