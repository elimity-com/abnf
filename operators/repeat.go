package operators

// Repeat defines a variable repetition.
func Repeat(key string, min, max int, r Operator) Operator {
	return func(s []byte) Alternatives {
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

// RepeatN defines a specific repetition.
func RepeatN(key string, n int, r Operator) Operator {
	return Repeat(key, n, n, r)
}

// Repeat1Inf defines a specific repetition from 0 to infinity.
func Repeat0Inf(key string, r Operator) Operator {
	return Repeat(key, 0, -1, r)
}

// Repeat0Inf defines a specific repetition from 1 to infinity.
func Repeat1Inf(key string, r Operator) Operator {
	return Repeat(key, 1, -1, r)
}

// RepeatOptional defines a specific repetition from 0 to 1. Behaves the same as Optional.
func RepeatOptional(key string, r Operator) Operator {
	return Repeat(key, 0, 1, r)
}
