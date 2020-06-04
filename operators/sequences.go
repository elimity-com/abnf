package operators

func Concat(key string, rules ...Operator) Operator {
	return func(s []rune) Alternatives {
		var concat func(l int, rules []Operator) Alternatives
		concat = func(l int, rules []Operator) Alternatives {
			if len(rules) == 0 {
				return Alternatives{
					{
						Key:   key,
						Value: s[:l],
					},
				}
			}

			var nodes Alternatives
			subNodes := rules[0](s[l:])
			for _, node := range subNodes {
				// add node as child of next nodes
				for _, n := range concat(l+len(node.Value), rules[1:]) {
					n.Children = append([]*Node{node}, n.Children...)
					nodes = append(nodes, n)
				}
			}
			return nodes
		}
		return concat(0, rules)
	}
}
func Alts(key string, rules ...Operator) Operator {
	return func(s []rune) Alternatives {
		var nodes Alternatives
		for _, rule := range rules {
			subNodes := rule(s)
			for _, node := range subNodes {
				nodes = append(nodes, &Node{
					Key:      key,
					Value:    node.Value,
					Children: Children{node},
				})
			}
		}
		return nodes
	}
}
