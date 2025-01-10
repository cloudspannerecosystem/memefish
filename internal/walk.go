package internal

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/cloudspannerecosystem/memefish/ast"
)

type Visitor interface {
	Visit(path []string, node ast.Node) Visitor
}

func Walk(node ast.Node, v Visitor) {
	walk(node, []string{"$"}, v)
}

type InspectFuncType = func(path []string, node ast.Node) bool
type inspectFuncVisitor InspectFuncType

func (i inspectFuncVisitor) Visit(path []string, node ast.Node) Visitor {
	if !i(path, node) {
		return nil
	}
	return i
}

func InspectSlice[T ast.Node](nodes []T, f InspectFuncType) {
	WalkSlice(nodes, inspectFuncVisitor(f))
}

func Inspect(node ast.Node, f InspectFuncType) {
	Walk(node, inspectFuncVisitor(f))
}

type VisitorFunc func(path []string, node ast.Node) Visitor

func (f VisitorFunc) Visit(path []string, node ast.Node) Visitor {
	return f(path, node)
}

func WalkSlice[T ast.Node](nodes []T, v Visitor) {
	for i, node := range nodes {
		walk(node, []string{fmt.Sprintf("$[%d]", i)}, v)
	}
}

func walk(node ast.Node, path []string, v Visitor) {
	v = v.Visit(path, node)
	if v == nil {
		return
	}

	val := reflect.ValueOf(node)
	for val.Kind() != reflect.Struct {
		switch val.Kind() {
		case reflect.Ptr, reflect.Interface:
			val = val.Elem()
		default:
			return
		}
	}

	for i := range val.NumField() {
		fieldPart := "." + val.Type().Field(i).Name
		val := val.Field(i)
		switch {
		case val.Kind() == reflect.Slice && val.Type().Elem().Implements(reflect.TypeFor[ast.Node]()):
			// for i, val := range val.Seq2() {
			for i := range val.Len() {
				walk(safeCast[ast.Node](val.Index(i)), slices.Concat(path, []string{fieldPart + "[" + fmt.Sprint(i) + "]"}), v)
			}
		case val.Type().Implements(reflect.TypeFor[ast.Node]()):
			walk(safeCast[ast.Node](val), slices.Concat(path, []string{fieldPart}), v)
		}
	}
}

// safeCast returns argument value as parameter type.
// It is used to avoid typed-nil.
func safeCast[T comparable](val reflect.Value) T {
	var zero T
	if val.IsNil() {
		return zero
	}

	v, ok := val.Interface().(T)
	if !ok {
		return zero
	}
	return v
}
