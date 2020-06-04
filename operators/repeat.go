package operators

func Repeat(key string, min, max int, r Operator) Operator {
	return func(s []rune) Alternatives {
		var repeat func(i, l int) Alternatives
		repeat = func(i, l int) Alternatives {
			var nodes Alternatives
			if max < 0 || i < max {
				subNodes := r(s[l:])
				for _, node := range subNodes {
					for _, n := range repeat(i+1, l+len(node.Value)) {
						n.Children = append([]*Node{node}, n.Children...)
						nodes = append(nodes, n)
					}
				}
			}
			if i < min {
				return nodes
			}
			node := Node{
				Key:   key,
				Value: s[:l],
			}
			return append(nodes, &node)
		}
		return repeat(0, 0)
	}
}

func RepeatN(key string, n int, r Operator) Operator {
	return Repeat(key, n, n, r)
}

func Repeat0Inf(key string, r Operator) Operator {
	return Repeat(key, 0, -1, r)
}

func Repeat1Inf(key string, r Operator) Operator {
	return Repeat(key, 1, -1, r)
}

func RepeatOptional(key string, r Operator) Operator {
	return Repeat(key, 0, 1, r)
}
