package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/k0kubun/pp"
	"github.com/MakeNowJust/memefish/pkg/parser"
)

func init() {
	flag.Parse()
}

func main() {
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

	log.SetOutput(os.Stderr)
	if flag.NArg() < 2 {
		log.Fatal("USAGE: go run ./tools/update-result [input path] [result path]")
	}

	inputPath := flag.Arg(0)
	resultPath := flag.Arg(1)
	log.Printf("input: %q, result: %q", inputPath, resultPath)

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

	inputs, err := ioutil.ReadDir(inputPath)
	if err != nil {
		log.Fatalf("error on reading input path: %v", err)
	}

	for _, in := range inputs {
		log.Printf("update: %q", in.Name())
		b, err := ioutil.ReadFile(filepath.Join(inputPath, in.Name()))
		if err != nil {
			log.Fatalf("error on reading input file: %v", err)
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

		err = ioutil.WriteFile(filepath.Join(resultPath, in.Name()+".txt"), buf.Bytes(), 0666)
		if err != nil {
			log.Fatalf("error on writing result file: %v", err)
		}
	}
}
