// Package memefish is the foundation to analyze Spanner SQL.
//
// NOTE: it is dummy package to provide the document root.
package memefish

import (
	_ "github.com/MakeNowJust/memefish/pkg/analyzer"
	_ "github.com/MakeNowJust/memefish/pkg/ast"
	_ "github.com/MakeNowJust/memefish/pkg/char"
	_ "github.com/MakeNowJust/memefish/pkg/parser"
	_ "github.com/MakeNowJust/memefish/pkg/token"
)
