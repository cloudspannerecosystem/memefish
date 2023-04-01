package analyzer_test

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cloudspannerecosystem/memefish/pkg/analyzer"
	"github.com/cloudspannerecosystem/memefish/pkg/char"
	"github.com/cloudspannerecosystem/memefish/pkg/parser"
	"github.com/cloudspannerecosystem/memefish/pkg/token"
	"gopkg.in/yaml.v2"
)

type TestData struct {
	SQL     string            `yaml:"SQL"`
	Results []*TestDataColumn `yaml:"Results"`
}

type TestDataColumn struct {
	Name string `yaml:"Name"`
	Type string `yaml:"Type"`
}

func TestAnalyzeQueryStatement(t *testing.T) {
	testdataPath := "./testdata/query"
	files, err := ioutil.ReadDir(testdataPath)
	if err != nil {
		t.Fatalf("error on reading testdata path: %v", err)
	}

	for _, file := range files {
		t.Run(file.Name(), func(t *testing.T) {
			bs, err := ioutil.ReadFile(filepath.Join(testdataPath, file.Name()))
			if err != nil {
				t.Fatalf("error on reading input file: %v", err)
			}

			var testdata *TestData
			err = yaml.Unmarshal(bs, &testdata)
			if err != nil {
				t.Fatalf("error on parsing YAML: %v", err)
			}

			file := &token.File{FilePath: file.Name() + ".SQL", Buffer: testdata.SQL}
			p := &parser.Parser{
				Lexer: &parser.Lexer{File: file},
			}
			stmt, err := p.ParseQuery()
			if err != nil {
				t.Fatalf("error on parsing SQL: %v", err)
			}

			a := &analyzer.Analyzer{
				File: file,
			}
			err = a.AnalyzeQueryStatement(stmt)
			if err != nil {
				t.Fatalf("error on analyzing SQL: %v", err)
			}

			list := a.NameLists[stmt.Query]
			if list == nil {
				t.Fatal("empty list")
			}

			if len(list) != len(testdata.Results) {
				t.Fatalf("result columns length is not matched: %d != %d", len(list), len(testdata.Results))
			}

			for i, name := range list {
				if name.Text != testdata.Results[i].Name {
					t.Fatalf("result column name is not matched, %q != %q", name.Text, testdata.Results[i].Name)
				}

				expected, ok := parseTypeTiny(testdata.Results[i].Type)
				if !ok {
					t.Fatalf("error on parsing expected column type: %s", testdata.Results[i].Type)
				}

				if !typeEqual(name.Type, expected) {
					t.Fatalf("result column type is not matched: %s != %s", analyzer.TypeString(name.Type), analyzer.TypeString(expected))
				}
			}
		})
	}
}

func parseTypeTiny(s string) (analyzer.Type, bool) {
	s = strings.TrimSpace(s)
	switch char.ToUpper(s) {
	case "NULL":
		return nil, true
	case "BOOL":
		return analyzer.BoolType, true
	case "INT64":
		return analyzer.Int64Type, true
	case "FLOAT64":
		return analyzer.Float64Type, true
	case "STRING":
		return analyzer.StringType, true
	case "BYTES":
		return analyzer.BytesType, true
	case "DATE":
		return analyzer.DateType, true
	case "TIMESTAMP":
		return analyzer.TimestampType, true
	case "NUMERIC":
		return analyzer.NumericType, true
	}

	if len(s) <= 6 {
		return nil, false
	}

	if char.EqualFold(s[:6], "ARRAY<") && s[len(s)-1] == '>' {
		t, ok := parseTypeTiny(s[6 : len(s)-1])
		if !ok {
			return nil, false
		}
		return &analyzer.ArrayType{Item: t}, true
	}

	if char.EqualFold(s[:7], "STRUCT<") && s[len(s)-1] == '>' {
		var fs []*analyzer.StructField
		i := 7
		for i < len(s) {
			j := i
			nest := 0
			for j < len(s) {
				if nest == 0 && (s[j] == '>' || s[j] == ',') {
					break
				}
				if s[j] == '<' {
					nest++
				}
				if s[j] == '>' {
					nest++
				}
				j++
			}
			ss := strings.SplitN(strings.TrimSpace(s[i:j]), " ", 2)
			switch len(ss) {
			case 0:
				return nil, false
			case 1:
				t, ok := parseTypeTiny(ss[0])
				if !ok {
					return nil, false
				}
				fs = append(fs, &analyzer.StructField{Type: t})
			case 2:
				t, ok := parseTypeTiny(ss[1])
				if !ok {
					return nil, false
				}
				fs = append(fs, &analyzer.StructField{Name: ss[0], Type: t})
			}
			i = j + 1
		}
		return &analyzer.StructType{Fields: fs}, true
	}

	return nil, false
}

func typeEqual(s, t analyzer.Type) bool {
	if s == t {
		return true
	}
	if s == nil || t == nil {
		return false
	}

	switch s := s.(type) {
	case analyzer.SimpleType:
		t, ok := t.(analyzer.SimpleType)
		if !ok {
			return false
		}
		return s == t
	case *analyzer.ArrayType:
		t, ok := t.(*analyzer.ArrayType)
		if !ok {
			return false
		}
		return typeEqual(s.Item, t.Item)
	case *analyzer.StructType:
		t, ok := t.(*analyzer.StructType)
		if !ok {
			return false
		}
		if len(s.Fields) != len(t.Fields) {
			return false
		}
		for i, f := range s.Fields {
			if f.Name != t.Fields[i].Name || !typeEqual(f.Type, t.Fields[i].Type) {
				return false
			}
		}
		return true
	}

	panic("BUG: unreachable")
}
