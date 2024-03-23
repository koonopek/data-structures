package main

import (
	"fmt"
	"strings"
)

const (
	NODE_ROOT = iota
	NODE_EXT
	NODE_LEAF
)

// when t == NODE_EXT then path is set and it behaves like edge
// when t == NODE_LEAF then value is set and it behaves like leaf
type Node struct {
	t      int
	value  []byte
	path   []byte
	childs []*Node
}

func (node *Node) Insert(key []byte, value []byte) {
	for i, _ := range node.childs {
		if node.childs[i].t == NODE_EXT {
			prefixBytesEqualCount := equalUpTo(key, node.childs[i].path)
			if prefixBytesEqualCount > 0 {

				newValueNode := Node{value: value, t: NODE_LEAF}
				newPath := node.childs[i].path[0:prefixBytesEqualCount]
				newExtNode := Node{path: key[prefixBytesEqualCount:], childs: []*Node{&newValueNode}, t: NODE_EXT}

				newCommonExtNode := Node{path: newPath, childs: []*Node{node.childs[i], &newExtNode}, t: NODE_EXT}

				node.childs[i].path = node.childs[i].path[prefixBytesEqualCount:]
				node.childs[i] = &newCommonExtNode
				return
			}

			if prefixBytesEqualCount == len(key) {
				panic(fmt.Sprintf("tried to insert duplicated key=%v", string(key)))
			}
		} else {
			panic("I should never be here")
		}
	}

	leafNode := Node{t: NODE_LEAF, value: value}
	node.childs = append(node.childs, &Node{t: NODE_EXT, path: key, childs: []*Node{&leafNode}})
}

func (node *Node) Find(key []byte) *Node {
	for _, child := range node.childs {
		prefixBytesEqualCount := equalUpTo(key, child.path)

		if prefixBytesEqualCount == len(key) {
			return child.childs[0]
		}

		if prefixBytesEqualCount > 0 {
			return child.Find(key[prefixBytesEqualCount:])
		}

	}

	return nil
}

func (node *Node) Print(depth int) {
	if len(node.childs) == 0 {
		return
	}

	for _, child := range node.childs {
		fmt.Printf("%v", strings.Repeat(" ", depth))
		if child.t == NODE_EXT {
			fmt.Printf("ext path=%v\n", string(child.path))
			depth++
			child.Print(depth)
		} else if child.t == NODE_LEAF {
			fmt.Printf("leaf value=%v\n", string(child.value))

		} else {
			panic("one of node doesnt have assigned type")
		}
	}
}

// 0 means that 0 first characters are equal
func equalUpTo(a []byte, b []byte) int {
	shorter := len(a)
	if shorter > len(b) {
		shorter = len(b)
	}

	for i := 0; i < shorter; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return shorter
}

func main() {
	root := Node{t: NODE_ROOT}

	root.Insert([]byte("ala"), []byte("kot"))
	root.Insert([]byte("alaw"), []byte("kot"))
	root.Insert([]byte("tomek"), []byte("pies"))
	root.Insert([]byte("t"), []byte("pies"))

	fmt.Println("Pritninting whole tree")
	root.Print(0)

	found := root.Find([]byte("tomek"))
	fmt.Printf("found %v\n", string(found.value))
}
