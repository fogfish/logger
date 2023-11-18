//
// Copyright (C) 2021 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/logger
//

package trie

import (
	"fmt"
	"log/slog"
	"strings"
)

// Node of trie
type Node struct {
	Path  string  // substring from the path "owned" by the node
	Heir  []*Node // heir nodes
	Level slog.Level
}

// New creates new trie
func New() *Node {
	root := &Node{
		Heir: []*Node{},
	}

	return root
}

// lookup is hot-path discovery of node at the path
func (root *Node) Lookup(path string) (at int, node *Node) {
	node = root
lookup:
	for {
		// leaf node, no futher lookup is possible
		// return current `node`` and position `at` path
		if len(node.Heir) == 0 {
			return
		}

		for _, heir := range node.Heir {
			if len(path[at:]) < len(heir.Path) {
				// No match, path cannot match node
				continue
			}

			// the node consumers entire path
			if len(heir.Path) == 1 && heir.Path[0] == '*' {
				at = len(path)
				node = heir
				return
			}

			if path[at] != heir.Path[0] {
				// No match, path cannot match node
				// this is micro-optimization to reduce overhead of memequal
				continue
			}

			if path[at:at+len(heir.Path)] == heir.Path {
				// node matches the path, continue lookup
				at = at + len(heir.Path)
				node = heir
				continue lookup
			}
		}

		return
	}
}

func (root *Node) Append(path string, level slog.Level) {
	if strings.HasSuffix(path, "*") {
		node := root.append(path[:len(path)-1], 0)
		node.Heir = append(node.Heir, &Node{Path: "*", Heir: make([]*Node, 0), Level: level})
		return
	}

	root.append(path, level)
}

func (root *Node) append(path string, level slog.Level) *Node {
	if len(path) == 0 {
		_, n := root.appendTo("/")
		n.Level = level
	}

	at, node := root.appendTo(path)

	// split the node and add endpoint
	if len(path[at:]) != 0 {
		split := &Node{
			Path: path[at:],
			Heir: make([]*Node, 0),
		}
		node.Heir = append(node.Heir, split)
		node = split
	}

	node.Level = level
	return node
}

// appendTo finds the node in trie where to add path (or segment).
// It returns the candidate node and length of "consumed" path
func (root *Node) appendTo(path string) (at int, node *Node) {
	node = root
lookup:
	for {
		if len(node.Heir) == 0 {
			// leaf node, no futher lookup is possible
			// return current `node`` and position `at` path
			return
		}

		for _, heir := range node.Heir {
			prefix := longestCommonPrefix(path[at:], heir.Path)
			at = at + prefix
			switch {
			case prefix == 0:
				// No common prefix, jump to next heir
				continue
			case prefix == len(heir.Path):
				// Common prefix is the node itself, continue lookup into heirs
				node = heir
				continue lookup
			default:
				// Common prefix is shorter than node itself, split is required
				if prefixNode := node.heirByPath(heir.Path[:prefix]); prefixNode != nil {
					// prefix already exists, current node needs to be moved
					// under existing one
					node.Path = node.Path[prefix:]
					prefixNode.Heir = append(prefixNode.Heir, node)
					node = prefixNode
					return
				}

				// prefix does not exist, current node needs to be split
				// the list of heirs needs to be patched
				for j := 0; j < len(node.Heir); j++ {
					if node.Heir[j].Path == heir.Path {
						n := heir
						node.Heir[j] = &Node{
							Path: heir.Path[:prefix],
							Heir: []*Node{n},
						}
						n.Path = heir.Path[prefix:]
						node = node.Heir[j]
						return
					}
				}
			}
		}
		// No heir is found return current node
		return
	}
}

func (root *Node) heirByPath(path string) *Node {
	for i := 0; i < len(root.Heir); i++ {
		if root.Heir[i].Path == path {
			return root.Heir[i]
		}
	}
	return nil
}

// Walk through trie, use for debug purposes only
func (root *Node) Walk(f func(int, *Node)) {
	walk(root, 0, f)
}

func walk(node *Node, level int, f func(int, *Node)) {
	f(level, node)
	for _, n := range node.Heir {
		walk(n, level+1, f)
	}
}

// Println outputs trie to console
func (root *Node) Println() {
	root.Walk(
		func(i int, n *Node) {
			fmt.Println(strings.Repeat(" ", i), n.Path)
		},
	)
}

//
// Utils
//

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func longestCommonPrefix(a, b string) (prefix int) {
	max := min(len(a), len(b))
	for prefix < max && a[prefix] == b[prefix] {
		prefix++
	}
	return
}
