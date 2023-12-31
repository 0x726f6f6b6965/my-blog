package utils

import (
	"sort"
	"sync"
)

type Node struct {
	children map[rune]*Node
	isEnd    bool
	data     string
	rank     int
	cache    []*Node
}

func (node *Node) Update(n *Node) {
	// find out if n node is in the cache
	nd := sort.Search(len(node.cache), func(i int) bool { return node.cache[i] == n })
	// if not append n node to the hot
	if nd <= len(node.cache) {
		node.cache = append(node.cache, n)
	}
	// sort with rank
	sort.Slice(node.cache, func(i, j int) bool {
		return node.cache[i].rank > node.cache[j].rank
	})
	if len(node.cache) > 5 {
		node.cache = node.cache[:len(node.cache)-1]
	}
}

type WordDictionary struct {
	root *Node
	sync.RWMutex
}

func NewWordDictionary() *WordDictionary {
	node := &Node{children: make(map[rune]*Node)}
	return &WordDictionary{root: node}
}

func (wordDict *WordDictionary) InsertWord(word string) {
	wordDict.RLock()
	defer wordDict.RUnlock()
	node := wordDict.root
	var visited []*Node
	for _, c := range word {
		if _, ok := node.children[c]; !ok {
			node.children[c] = &Node{children: make(map[rune]*Node)}
		}
		node = node.children[c]
		visited = append(visited, node)
	}
	node.isEnd = true
	node.data = word
	node.rank += 1
	for _, n := range visited {
		n.Update(node)
	}
}

func (wordDict *WordDictionary) SearchWord(word string) []string {
	node := wordDict.root
	var res []string
	for _, c := range word {
		if _, ok := node.children[c]; !ok {
			break
		}
		node = node.children[c]
	}

	for _, node := range node.cache {
		res = append(res, node.data)
	}
	return res
}
