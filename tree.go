package abnf

type AST struct {
	Key      string
	Value    []rune
	Children []AST
}

func (ast *AST) String() string {
	return string(ast.Value)
}

func (ast *AST) ValueString() string {
	return string(ast.Value)
}

func (ast *AST) GetNode(name string, recursive bool) *AST {
	for _, child := range ast.Children {
		if child.Key == name {
			return &child
		}
		if recursive {
			if n := child.GetNode(name, recursive); n != nil {
				return n
			}
		}
	}
	return nil
}

func (ast *AST) GetAllNodes(name string) []AST {
	nodes := make([]AST, 0)
	for _, child := range ast.Children {
		if child.Key == name {
			nodes = append(nodes, child)
		}
		nodes = append(nodes, child.GetAllNodes(name)...)
	}
	return nodes
}
