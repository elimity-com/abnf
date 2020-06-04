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

func (n *Node) GetImmediateSubNode(key string) *Node {
	return n.Children.Get(key, false)
}

func (n *Node) GetSubNode(key string) *Node {
	return n.Children.Get(key, true)
}

func (n *Node) GetSubNodes(key string) Alternatives {
	return n.Children.GetAll(key)
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

func (c Children) GetAll(key string) Alternatives {
	var nodes Alternatives
	for _, child := range c {
		if child.Key == key {
			nodes = append(nodes, child)
		}
		nodes = append(nodes, child.GetSubNodes(key)...)
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

type Operator func([]rune) Alternatives
