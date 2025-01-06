package astcatalog

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
)

var (
	rePosLine   = regexp.MustCompile(`(?m)^\s*pos\s*=\s*(.*)`)
	reEndLine   = regexp.MustCompile(`(?m)^\s*end\s*=\s*(.*)`)
	reTmplLines = regexp.MustCompile(`(?m)((?:^\t.*\n)+)+`)
)

func Load(filename string) (Catalog, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	catalog := make(Catalog)
	interfaces := make(map[NodeInterfaceType]struct{})

	commentMap := ast.NewCommentMap(fset, f, f.Comments)

	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					switch t := s.Type.(type) {
					case *ast.StructType:
						name := NodeStructType(s.Name.Name)
						if _, ok := catalog[name]; !ok {
							catalog[name] = &NodeDef{
								Name: s.Name.Name,
							}
						}

						node := catalog[name]
						node.Doc = d.Doc.Text()
						node.SourcePos = s.Pos()

						if m := reTmplLines.FindAllStringSubmatch(node.Doc, -1); m != nil {
							node.Tmpl = m[len(m)-1][0]
						} else {
							return nil, fmt.Errorf("no template found: %s", name)
						}

						comments := commentMap.Filter(t).Comments()
						if len(comments) == 0 {
							return nil, fmt.Errorf("no pos/end comment found: %s", name)
						}

						// We assume the first comment group should a pos/end comment.
						posComment := comments[0].Text()
						if m := rePosLine.FindStringSubmatch(posComment); m != nil {
							node.Pos = m[1]
						} else {
							return nil, fmt.Errorf("no pos comment found: %s", name)
						}
						if m := reEndLine.FindStringSubmatch(posComment); m != nil {
							node.End = m[1]
						} else {
							return nil, fmt.Errorf("no end coment found: %s", name)
						}

						for _, f := range t.Fields.List {
							ft, err := loadType(f.Type, interfaces)
							if err != nil {
								return nil, err
							}

							comment := ""
							comments := commentMap.Filter(f).Comments()
							for _, c := range comments {
								if f.End() < c.Pos() {
									comment = c.Text()
									break
								}
							}
							for _, name := range f.Names {
								node.Fields = append(node.Fields, &FieldDef{
									Name:    name.Name,
									Type:    ft,
									Comment: comment,
								})
							}
						}
					case *ast.InterfaceType:
						name := NodeInterfaceType(s.Name.Name)
						if _, ok := interfaces[name]; ok {
							return nil, fmt.Errorf("duplicated interface: %s", name)
						}

						interfaces[name] = struct{}{}
					default:
						return nil, fmt.Errorf("unexpected spec: %#v", t)
					}
				}
			}
		case *ast.FuncDecl:
			if d.Recv == nil {
				return nil, fmt.Errorf("unexpected func decl: %#v", d)
			}

			recv, err := loadType(d.Recv.List[0].Type, interfaces)
			if err != nil {
				return nil, err
			}

			structName, ok := recv.(NodeStructType)
			if !ok {
				return nil, fmt.Errorf("unexpected receiver type: %#v", recv)
			}

			funcName := d.Name.Name
			cutName, found := strings.CutPrefix(funcName, "is")
			if !found {
				return nil, fmt.Errorf("unexpected func name: %s", funcName)
			}

			interfaceName := NodeInterfaceType(cutName)
			if _, ok := interfaces[interfaceName]; !ok {
				return nil, fmt.Errorf("unknown interface: %s", interfaceName)
			}

			if _, ok := catalog[structName]; !ok {
				catalog[structName] = &NodeDef{
					Name: string(structName),
				}
			}

			node := catalog[structName]
			node.Implements = append(node.Implements, interfaceName)
		}
	}

	return catalog, nil
}

func loadType(t ast.Expr, interfaces map[NodeInterfaceType]struct{}) (FieldType, error) {
	switch t := t.(type) {
	case *ast.Ident:
		switch t.Name {
		case "bool":
			return BoolType, nil
		case "int":
			return IntType, nil
		case "string":
			return StringType, nil
		default:
			if _, ok := interfaces[NodeInterfaceType(t.Name)]; ok {
				return NodeInterfaceType(t.Name), nil
			}
			return NodeStructType(t.Name), nil
		}
	case *ast.SelectorExpr:
		if x, ok := t.X.(*ast.Ident); !(ok && x.Name == "token") {
			return nil, fmt.Errorf("unexpected selector expr: %#v", t)
		}
		switch t.Sel.Name {
		case "Pos":
			return TokenPosType, nil
		case "Token":
			return TokenTokenType, nil
		default:
			return nil, fmt.Errorf("unexpected selector name: %#v", t)
		}
	case *ast.StarExpr:
		ft, err := loadType(t.X, interfaces)
		if err != nil {
			return nil, err
		}

		return PointerType{Type: ft}, nil
	case *ast.ArrayType:
		if t.Len != nil {
			return nil, fmt.Errorf("unexpected array type: %#v", t)
		}

		ft, err := loadType(t.Elt, interfaces)
		if err != nil {
			return nil, err
		}

		return SliceType{Type: ft}, nil
	default:
		return nil, fmt.Errorf("unexpected type: %#v", t)
	}
}
