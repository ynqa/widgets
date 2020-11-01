package node

import (
	"fmt"
	"sort"
	"strings"
)

type Node struct {
	name         string
	row          []string
	isRoot       bool
	isLeaf       bool
	depth        int
	childVisible bool
	parent       *Node
	children     []*Node
}

func Root() *Node {
	return &Node{
		isRoot:       true,
		childVisible: true,
		children:     make([]*Node, 0),
	}
}

func New(name string, row []string) *Node {
	return &Node{
		name:     name,
		row:      row,
		children: make([]*Node, 0),
	}
}

func Leaf(name string, row []string) *Node {
	return &Node{
		name:   name,
		row:    row,
		isLeaf: true,
	}
}

func (n *Node) String() string {
	nodes := n.FlattenAll()
	var builder strings.Builder
	for _, n := range nodes {
		builder.WriteString(fmt.Sprintf("%s%s: %v\n", strings.Repeat(" ", n.Depth()), n.name, n.row))
	}
	return builder.String()
}

func (n *Node) ChildVisible() bool {
	return n.childVisible
}

func (n *Node) Depth() int {
	return n.depth
}

func (n *Node) IsLeaf() bool {
	return n.isLeaf
}

func (n *Node) Row() []string {
	return n.row
}

func (n *Node) Append(children ...*Node) *Node {
	if len(children) > 0 && !n.isLeaf {
		for _, child := range children {
			child.parent = n
		}
		n.children = append(n.children, children...)
		sort.Slice(n.children, func(i, j int) bool {
			return n.children[i].name < n.children[j].name
		})
	}
	return n
}

func ApplyChildVisible(old *Node, new *Node) *Node {
	if old == nil && new == nil {
		return nil
	}
	if old == nil {
		return new
	}
	if new == nil {
		return old
	}

	if new.name == old.name {
		new.childVisible = old.childVisible
	}

	for _, c1 := range old.children {
		for _, c2 := range new.children {
			ApplyChildVisible(c1, c2)
		}
	}
	return new
}

func (n *Node) Toggle(idx int) {
	nodes := n.Flatten()
	if 0 <= idx && idx < len(nodes) {
		nodes[idx].childVisible = !nodes[idx].childVisible
	}
}

func (n *Node) Flatten() []*Node {
	return n.flatten(-1, false)
}

func (n *Node) FlattenAll() []*Node {
	return n.flatten(-1, true)
}

func (n *Node) flatten(cursor int, all bool) []*Node {
	var nodes []*Node
	if n != nil {
		if !n.isRoot {
			n.depth = cursor
			nodes = append(nodes, n)
		}
		if all || n.childVisible {
			for _, node := range n.children {
				nodes = append(nodes, node.flatten(cursor+1, all)...)
			}
		}
	}
	return nodes
}

func (n *Node) Names() []string {
	var names []string
	if !n.isRoot {
		names = append(names, n.name)
	}
	if n.parent != nil {
		names = append(names, n.parent.Names()...)
	}
	return names
}
