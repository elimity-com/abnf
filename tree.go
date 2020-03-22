package abnf

type AST struct {
	Name     string
	Raw      []rune
	Children []AST
}

func (ast *AST) String() string {
	return string(ast.Raw)
}

func (ast *AST) GetNode(name string) *AST {
	for _, child := range ast.Children {
		if child.Name == name {
			return &child
		}
	}
	return nil
}

func (ast *AST) GetAllNodes(name string) []AST {
	nodes := make([]AST, 0)
	for _, child := range ast.Children {
		if child.Name == name {
			nodes = append(nodes, child)
		}
		nodes = append(nodes, child.GetAllNodes(name)...)
	}
	return nodes
}
