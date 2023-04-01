// Package memefish is the foundation to analyze Spanner SQL.
//
// NOTE: it is dummy package to provide the document root.
package memefish

import (
	_ "github.com/cloudspannerecosystem/memefish/pkg/analyzer"
	_ "github.com/cloudspannerecosystem/memefish/pkg/ast"
	_ "github.com/cloudspannerecosystem/memefish/pkg/char"
	_ "github.com/cloudspannerecosystem/memefish/pkg/parser"
	_ "github.com/cloudspannerecosystem/memefish/pkg/token"
)
