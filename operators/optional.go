package operators

func Optional(key string, r Operator) Operator {
	return func(s []rune) Alternatives {
		empty := &Node{
			Key:   key,
			Value: s[:0],
		}

		subNodes := r(s)
		if len(subNodes) == 0 {
			return Alternatives{empty}
		}

		var nodes Alternatives
		for _, node := range subNodes {
			nodes = append(nodes, &Node{
				Key:      key,
				Value:    node.Value,
				Children: Children{node},
			})
		}
		return append(nodes, empty)
	}
}
