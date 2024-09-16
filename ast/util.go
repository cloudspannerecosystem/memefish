package ast

// sqlOpt output:
//
//	when node != nil: left + node.SQL() + right
//	when node == nil: empty string
//
// requires Go 1.20
func sqlOpt[T interface {
	Node
	comparable
}](left string, node T, right string) string {
	var zero T
	if node == zero {
		return ""
	}
	return left + node.SQL() + right
}
