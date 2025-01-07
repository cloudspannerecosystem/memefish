package astcatalog

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
)

// Load loads a catalog from the given AST files.
func Load(astFilename, astConstFilename string) (*Catalog, error) {
	fset := token.NewFileSet()

	consts, err := loadConsts(fset, astConstFilename)
	if err != nil {
		return nil, err
	}

	structs, interfaces, err := loadStructs(fset, astFilename, consts)
	if err != nil {
		return nil, err
	}

	catalog := &Catalog{
		Structs:    structs,
		Interfaces: interfaces,
		Consts:     consts,
	}
	return catalog, nil
}

func loadConsts(fset *token.FileSet, filename string) (map[ConstType]*ConstDef, error) {
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse '%s': %w", filename, err)
	}

	consts := make(map[ConstType]*ConstDef)

	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					name := ConstType(s.Name.Name)
					if _, ok := consts[name]; ok {
						return nil, fmt.Errorf("duplicated const: %s", name)
					}

					consts[name] = &ConstDef{
						SourcePos: d.Pos(),
						Name:      s.Name.Name,
					}
				case *ast.ValueSpec:
					if s.Type == nil {
						return nil, fmt.Errorf("unexpected value spec: %#v", s)
					}
					ty, err := loadType(s.Type, nil, consts)
					if err != nil {
						return nil, err
					}
					name, ok := ty.(ConstType)
					if !ok {
						return nil, fmt.Errorf("unexpected type: %#v", s.Type)
					}

					constDef, ok := consts[name]
					if !ok {
						return nil, fmt.Errorf("unknown const: %s", name)
					}

					if len(s.Values) != 1 {
						return nil, fmt.Errorf("unexpected values: %#v", s.Values)
					}
					lit, ok := s.Values[0].(*ast.BasicLit)
					if !(ok && lit.Kind == token.STRING) {
						return nil, fmt.Errorf("unexpected value: %#v", s.Values[0])
					}
					v := strings.Trim(lit.Value, "\"")

					for _, name := range s.Names {
						constDef.Values = append(constDef.Values, &ConstValueDef{
							Name:  name.Name,
							Value: v,
						})
					}
				}
			}
		}
	}

	return consts, nil
}

// Regular expressions to extract pos/end and template comments.
var (
	rePosLine   = regexp.MustCompile(`(?m)^\s*pos\s*=\s*(.*)`)
	reEndLine   = regexp.MustCompile(`(?m)^\s*end\s*=\s*(.*)`)
	reTmplLines = regexp.MustCompile(`(?m)((?:^\t.*\n)+)+`)
)

func loadStructs(fset *token.FileSet, filename string, consts map[ConstType]*ConstDef) (map[NodeStructType]*NodeStructDef, map[NodeInterfaceType]*NodeInterfaceDef, error) {
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse '%s': %w", filename, err)
	}

	structs := make(map[NodeStructType]*NodeStructDef)
	interfaces := make(map[NodeInterfaceType]*NodeInterfaceDef)

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
						if _, ok := structs[name]; !ok {
							structs[name] = &NodeStructDef{
								Name: s.Name.Name,
							}
						}

						structDef := structs[name]
						structDef.SourcePos = s.Pos()
						structDef.Doc = d.Doc.Text()

						if m := reTmplLines.FindAllStringSubmatch(structDef.Doc, -1); m != nil {
							structDef.Tmpl = m[len(m)-1][0]
						} else {
							return nil, nil, fmt.Errorf("no template found: %s", name)
						}

						comments := commentMap.Filter(t).Comments()
						if len(comments) == 0 {
							return nil, nil, fmt.Errorf("no pos/end comment found: %s", name)
						}

						// We assume the first comment group in the struct should a pos/end comment.
						posComment := comments[0].Text()
						if m := rePosLine.FindStringSubmatch(posComment); m != nil {
							structDef.Pos = m[1]
						} else {
							return nil, nil, fmt.Errorf("no pos comment found: %s", name)
						}
						if m := reEndLine.FindStringSubmatch(posComment); m != nil {
							structDef.End = m[1]
						} else {
							return nil, nil, fmt.Errorf("no end coment found: %s", name)
						}

						for _, f := range t.Fields.List {
							ty, err := loadType(f.Type, interfaces, consts)
							if err != nil {
								return nil, nil, err
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
								structDef.Fields = append(structDef.Fields, &FieldDef{
									Name:    name.Name,
									Type:    ty,
									Comment: comment,
								})
							}
						}
					case *ast.InterfaceType:
						name := NodeInterfaceType(s.Name.Name)
						if _, ok := interfaces[name]; ok {
							return nil, nil, fmt.Errorf("duplicated interface: %s", name)
						}

						// Node interface is special, so we skip it.
						if name == "Node" {
							continue
						}

						interfaces[name] = &NodeInterfaceDef{
							SourcePos: s.Pos(),
							Name:      string(name),
						}
					default:
						return nil, nil, fmt.Errorf("unexpected spec: %#v", t)
					}
				}
			}

		case *ast.FuncDecl:
			if d.Recv == nil || len(d.Recv.List) != 1 {
				return nil, nil, fmt.Errorf("unexpected func decl: %#v", d)
			}

			recv, err := loadType(d.Recv.List[0].Type, interfaces, consts)
			if err != nil {
				return nil, nil, err
			}

			structName, ok := recv.(NodeStructType)
			if !ok {
				return nil, nil, fmt.Errorf("unexpected receiver type: %#v", recv)
			}

			funcName := d.Name.Name
			cutName, found := strings.CutPrefix(funcName, "is")
			if !found {
				return nil, nil, fmt.Errorf("unexpected func name: %s", funcName)
			}

			interfaceName := NodeInterfaceType(cutName)
			interfaceDef, ok := interfaces[interfaceName]
			if !ok {
				return nil, nil, fmt.Errorf("unknown interface: %s", interfaceName)
			}

			interfaceDef.Implemented = append(interfaceDef.Implemented, structName)

			if _, ok := structs[structName]; !ok {
				structs[structName] = &NodeStructDef{
					Name: string(structName),
				}
			}

			structDef := structs[structName]
			structDef.Implements = append(structDef.Implements, interfaceName)
		}
	}

	return structs, interfaces, nil
}

func loadType(t ast.Expr, interfaces map[NodeInterfaceType]*NodeInterfaceDef, consts map[ConstType]*ConstDef) (Type, error) {
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
			if _, ok := consts[ConstType(t.Name)]; ok {
				return ConstType(t.Name), nil
			}
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
		ty, err := loadType(t.X, interfaces, consts)
		if err != nil {
			return nil, err
		}

		return PointerType{Type: ty}, nil
	case *ast.ArrayType:
		if t.Len != nil {
			return nil, fmt.Errorf("unexpected array type: %#v", t)
		}

		ty, err := loadType(t.Elt, interfaces, consts)
		if err != nil {
			return nil, err
		}

		return SliceType{Type: ty}, nil
	default:
		return nil, fmt.Errorf("unexpected type: %#v", t)
	}
}
