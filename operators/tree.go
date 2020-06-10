package operators

import (
	"fmt"
	"strings"
)

type Node struct {
	Key      string
	Value    []rune
	Children Children
}

func (n *Node) String() string {
	return string(n.Value)
}

func (n *Node) StringRecursive() string {
	return n.stringRecursive(1)
}

func (n *Node) stringRecursive(depth int) string {
	str := fmt.Sprintf("%s %s: %s\n", strings.Repeat("-", depth), n.Key, string(n.Value))
	for _, child := range n.Children {
		str += child.stringRecursive(depth + 1)
	}
	return str
}

func (n *Node) Equals(other *Node) error {
	if n == nil || other == nil {
		return fmt.Errorf("one of the nodes is nil")
	}

	if n.Key != other.Key {
		return fmt.Errorf("keys do not match: %s, %s", n.Key, other.Key)
	}
	if len(n.Value) != len(other.Value) {
		return fmt.Errorf("values do not match: %s, %s", string(n.Value), string(other.Value))
	}
	if len(n.Children) != len(other.Children) {
		return fmt.Errorf("not an equal amount of children: %d, %d", len(n.Children), len(other.Children))
	}

	for i := range n.Value {
		if n.Value[i] != other.Value[i] {
			return fmt.Errorf("value does not match, index %d", i)
		}
	}

	for i := range n.Children {
		if err := n.Children[i].Equals(other.Children[i]); err != nil {
			return err
		}
	}

	return nil
}

func (n *Node) IsEmpty() bool  {
	return len(n.Value) == 0
}

func (n *Node) GetImmediateSubNode(key string) *Node {
	return n.Children.Get(key, false)
}

func (n *Node) GetNode(key string) *Node {
	if n.Key == key {
		return n
	}
	return n.GetSubNode(key)
}

func (n *Node) GetSubNode(key string) *Node {
	return n.Children.Get(key, true)
}

func (n *Node) GetSubNodes(key string) Children {
	return n.Children.GetAll(key)
}

func (n *Node) GetSubNodesBefore( key, before string) Children {
	return n.Children.GetAllBefore(key, before)
}

func (n *Node) Contains(key string) bool {
	for _, child := range n.Children {
		if child.Key == key ||
			child.Contains(key) {
			return true
		}
	}
	return false
}

type Children []*Node

func (c Children) Get(key string, recursive bool) *Node {
	for _, child := range c {
		if child.Key == key {
			return child
		}
		if recursive {
			if n := child.GetSubNode(key); n != nil {
				return n
			}
		}
	}
	return nil
}

func (c Children) GetAll(key string) Children {
	var nodes Children
	for _, child := range c {
		if child.Key == key {
			nodes = append(nodes, child)
		}
		nodes = append(nodes, child.GetSubNodes(key)...)
	}
	return nodes
}

func (c Children) GetAllBefore(key, before string) Children {
	var nodes Children
	for _, child := range c {
		if child.Key == before {
			return nodes
		}

		if child.Key == key {
			nodes = append(nodes, child)
		}
		nodes = append(nodes, child.GetSubNodesBefore(key, before)...)
	}
	return nodes
}

type Alternatives []*Node

func (as Alternatives) Best() *Node {
	best := &Node{}
	for _, a := range as {
		if len(a.Value) > len(best.Value) {
			best = a
		}
	}
	return best
}

func (as Alternatives) Equals(other Alternatives) error {
	if len(as) != len(other) {
		return fmt.Errorf("not an equal amount of nodes: %d, %d", len(as), len(other))
	}

	for i, alternative := range as {
		if err := alternative.Equals(other[i]); err != nil {
			return err
		}
	}
	return nil
}

type Operator func([]rune) Alternatives
