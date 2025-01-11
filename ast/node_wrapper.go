package ast

type comparableNode = interface {
	comparable
	Node
}

func wrapNode[T comparableNode](node T) Node {
	var zero T
	if node == zero {
		return nil
	}
	return node
}

func wrapNodes[T Node](nodes []T) []Node {
	result := make([]Node, 0, len(nodes))
	for _, node := range nodes {
		result = append(result, node)
	}
	return result
}
