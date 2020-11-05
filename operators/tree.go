package operators

import (
	"fmt"
	"regexp"
	"strings"
)

// Node represents a single node in a tree.
type Node struct {
	Key      string
	Value    []byte
	Children Children
}

// String returns the string representation of the node without new lines and duplicate spaces.
func (n *Node) String() string {
	return regexp.MustCompile(`\s+`).ReplaceAllString(string(n.Value), " ")
}

// StringRecursive returns the string representation of the whole (sub)tree/
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

// Equals checks whether the node's values is equal to another's.
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

// IsEmpty returns whether the node has no value.
func (n *Node) IsEmpty() bool {
	return len(n.Value) == 0
}

// GetImmediateSubNode searches the children of the node for the given key and returns that child.
func (n *Node) GetImmediateSubNode(key string) *Node {
	return n.Children.Get(key, false)
}

// GetNode searches itself and ALL the children of the node for the given key and returns that child.
func (n *Node) GetNode(key string) *Node {
	if n.Key == key {
		return n
	}
	return n.GetSubNode(key)
}

// GetSubNode searches ALL the children of the node for the given key and returns that child.
func (n *Node) GetSubNode(key string) *Node {
	return n.Children.Get(key, true)
}

// GetSubNode searches ALL the children of the node for the given key and returns all matching children.
func (n *Node) GetSubNodes(key string) Children {
	return n.Children.GetAll(key)
}

// GetSubNode searches ALL the children of the node for the given key and returns that child without passing a stop-key.
func (n *Node) GetSubNodesBefore(key string, stop ...string) Children {
	return n.Children.GetAllBefore(key, stop...)
}

// Contains returns whether the subtree contains the given key.
func (n *Node) Contains(key string) bool {
	for _, child := range n.Children {
		if child.Key == key ||
			child.Contains(key) {
			return true
		}
	}
	return false
}

// Children represent the children of a node.
type Children []*Node

// Get searches the children (recursively) for the given key and returns that child.
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

// GetAll returns ALL the children matching the given key.
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

// GetSubNode returns ALL the children matching the given key without passing a stop-key.
func (c Children) GetAllBefore(key string, before ...string) Children {
	var nodes Children
	for _, child := range c {
		for _, b := range before {
			if child.Key == b {
				return nodes
			}
		}

		if child.Key == key {
			nodes = append(nodes, child)
		}
		nodes = append(nodes, child.GetSubNodesBefore(key, before...)...)
	}
	return nodes
}

// Alternatives represent all alternative solutions for a tree. (Also contains subtrees, etc.)
type Alternatives []*Node

// Best returns the node with the best (longest) value.
func (as Alternatives) Best() *Node {
	best := &Node{}
	for _, a := range as {
		if len(a.Value) > len(best.Value) {
			best = a
		}
	}
	return best
}

// Equals checks whether the nodes are is equal to the given others.
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

// Operator represent an ABNF operator.
type Operator func([]byte) Alternatives
