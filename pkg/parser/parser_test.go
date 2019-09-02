package parser

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"github.com/k0kubun/pp"
	"github.com/pmezard/go-difflib/difflib"
)

func TestParser(t *testing.T) {
	// Disable color output.
	// https://github.com/k0kubun/pp/issues/26
	printer := pp.New()
	printer.SetColorScheme(pp.ColorScheme{
		Bool:            pp.NoColor,
		Integer:         pp.NoColor,
		Float:           pp.NoColor,
		String:          pp.NoColor,
		StringQuotation: pp.NoColor,
		EscapedChar:     pp.NoColor,
		FieldName:       pp.NoColor,
		PointerAdress:   pp.NoColor,
		Nil:             pp.NoColor,
		Time:            pp.NoColor,
		StructName:      pp.NoColor,
		ObjectLength:    pp.NoColor,
	})

	inputPath := "./testdata/input/query"
	resultPath := "./testdata/result/query"

	inputs, err := ioutil.ReadDir(inputPath)
	if err != nil {
		t.Fatalf("error on reading input path: %v", err)
	}

	for _, in := range inputs {
		t.Run(in.Name(), func(t *testing.T) {
			b, err := ioutil.ReadFile(filepath.Join(inputPath, in.Name()))
			if err != nil {
				t.Fatalf("error on reading input file: %v", err)
			}
			expected, err := ioutil.ReadFile(filepath.Join(resultPath, in.Name()+".txt"))
			if err != nil {
				t.Fatalf("error on reading result file: %v", err)
			}

			p := &Parser{
				Lexer: &Lexer{
					File: NewFile(in.Name(), string(b)),
				},
			}

			stmt, err := p.ParseQuery()
			if err != nil {
				log.Fatalf("error on parsing input file: %v", err)
			}

			var buf bytes.Buffer

			fmt.Fprintf(&buf, "--- %s\n", in.Name())
			fmt.Fprint(&buf, string(b))
			fmt.Fprintln(&buf)

			fmt.Fprintf(&buf, "--- AST\n")
			_, _ = printer.Fprintln(&buf, stmt)
			fmt.Fprintln(&buf)

			fmt.Fprintf(&buf, "--- SQL\n")
			fmt.Fprintln(&buf, stmt.SQL())

			actual := buf.Bytes()

			if bytes.Equal(actual, expected) {
				return
			}

			diff := difflib.UnifiedDiff{
				A:       difflib.SplitLines(string(expected)),
				B:       difflib.SplitLines(string(actual)),
				Context: 5,
			}
			d, err := difflib.GetUnifiedDiffString(diff)
			if err != nil {
				t.Fatalf("error on diff: %v", err)
			}
			t.Error(d)
		})
	}
}
