package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/MakeNowJust/memefish/pkg/analyzer"
	"github.com/MakeNowJust/memefish/pkg/parser"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
)

var param = flag.String("param", "", "param file")

func init() {
	flag.Parse()
}

func main() {
	if flag.NArg() < 1 {
		log.Fatal("usage: ./analyze [SQL query]")
	}

	query := flag.Arg(0)

	var params map[string]interface{}
	if *param != "" {
		log.Printf("load param file: %s", *param)
		var err error
		params, err = loadParamFile(*param)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("query: %q", query)

	p := &parser.Parser{
		Lexer: &parser.Lexer{
			File: parser.NewFile("[query]", query),
		},
	}

	log.Printf("start parsing")
	stmt, err := p.ParseQuery()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("finish parsing successfully")

	log.Printf("start analyzing")
	a := &analyzer.Analyzer{
		File:   p.File,
		Params: params,
	}
	a.AnalyzeQueryStatement(stmt)
	log.Printf("finish analyzing")

	list := a.NameLists[stmt.Query]
	if list == nil {
		log.Fatal("missing select list")
	}

	pp.Println(list)

	table := tablewriter.NewWriter(os.Stdout)
	var header []string
	for _, name := range list {
		header = append(header, name.Quote())
	}
	table.SetHeader(header)

	var types []string
	for _, name := range list {
		types = append(types, analyzer.TypeString(name.Type))
	}
	table.Append(types)

	table.Render()
}

type Param struct {
	BOOL    *bool               `yaml:"BOOL,omitempty"`
	INT64   *int64              `yaml:"INT64,omitempty"`
	FLOAT64 *float64            `yaml:"FLOAT64,omitempty"`
	STRING  *string             `yaml:"STRING,omitempty"`
	ARRAY   []*Param            `yaml:"ARRAY,omitempty"`
	STRUCT  []map[string]*Param `yaml:"STRUCT,omitempty"`
}

func loadParamFile(file string) (map[string]interface{}, error) {
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var params map[string]*Param
	err = yaml.Unmarshal(bs, &params)
	if err != nil {
		return nil, err
	}

	normalized := make(map[string]interface{})
	for name, p := range params {
		normalized[strings.ToUpper(name)] = decodeParam(p)
	}
	return normalized, nil
}

func decodeParam(p *Param) interface{} {
	switch {
	case p.BOOL != nil:
		return *p.BOOL
	case p.INT64 != nil:
		return *p.INT64
	case p.FLOAT64 != nil:
		return *p.FLOAT64
	case p.STRING != nil:
		return *p.STRING
	case p.ARRAY != nil:
		var result []interface{}
		for _, v := range p.ARRAY {
			result = append(result, decodeParam(v))
		}
		return result
	case p.STRUCT != nil:
		var result []map[string]interface{}
		for _, kv := range p.STRUCT {
			kvs := make(map[string]interface{})
			for name, v := range kv {
				kvs[name] = decodeParam(v)
			}
			result = append(result, kvs)
		}
		return result
	}

	panic("invalid param")
}
