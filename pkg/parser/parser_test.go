package parser_test

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/k0kubun/pp"
	"github.com/MakeNowJust/memefish/pkg/parser"
	"github.com/pmezard/go-difflib/difflib"
)

var update = flag.Bool("update", false, "update result files")

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

	if *update {
		_, err := os.Stat(resultPath)
		if err == nil {
			err = os.RemoveAll(resultPath)
			if err != nil {
				log.Fatalf("error on remove result path: %v", err)
			}
		}
		err = os.MkdirAll(resultPath, 0777)
		if err != nil {
			log.Fatalf("error on create result path: %v", err)
		}
	}

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

			p := &parser.Parser{
				Lexer: &parser.Lexer{
					File: parser.NewFile(in.Name(), string(b)),
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

			if *update {
				t.Log("update " + in.Name() + ".txt")
				err = ioutil.WriteFile(filepath.Join(resultPath, in.Name()+".txt"), buf.Bytes(), 0666)
				if err != nil {
					log.Fatalf("error on writing result file: %v", err)
				}
				return
			}

			expected, err := ioutil.ReadFile(filepath.Join(resultPath, in.Name()+".txt"))
			if err != nil {
				t.Fatalf("error on reading result file: %v", err)
			}

			if !bytes.Equal(actual, expected) {
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
				return
			}

			s1 := stmt.SQL()
			p1 := &parser.Parser{
				Lexer: &parser.Lexer{
					File: parser.NewFile(in.Name()+"(SQL)", s1),
				},
			}

			stmt1, err := p1.ParseQuery()
			if err != nil {
				log.Fatalf("error on parsing unparsed SQL: %v", err)
			}

			s2 := stmt1.SQL()
			if s1 != s2 {
				t.Errorf("%q != %q", s1, s2)
			}
		})
	}
}
