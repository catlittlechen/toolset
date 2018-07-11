// Author: catlittlechen@gmail.com

package lru

import (
	"errors"
	"fmt"
	"sync"
)

type Tree interface {
	Get(string) *Node
	Set(string, *Node)
	Del(string) *Node
}

type DefaultTree struct {
	tree map[string]*Node
}

func (dt *DefaultTree) Get(key string) *Node {
	return dt.tree[key]
}

func (dt *DefaultTree) Set(key string, node *Node) {
	dt.tree[key] = node
	return
}

func (dt *DefaultTree) Del(key string) *Node {
	node := dt.tree[key]
	delete(dt.tree, key)
	return node
}

type Node struct {
	pre  *Node
	next *Node

	key   string
	value interface{}
}

func (node *Node) clear() {
	node.pre = nil
	node.next = nil
	node.key = ""
	node.value = nil
	return
}

func (node *Node) init(k string, v interface{}) {
	node.key = k
	node.value = v
	return
}

type LRU struct {
	l *sync.Mutex

	left int

	head *Node
	tail *Node

	leftNode *Node
	tree     Tree
}

func New(size int, tree Tree) (*LRU, error) {
	if size <= 0 {
		return nil, errors.New("size must be large than 0")
	}

	if tree == nil {
		tree = &DefaultTree{
			tree: make(map[string]*Node),
		}
	}

	lru := &LRU{
		l:        new(sync.Mutex),
		left:     size,
		head:     nil,
		tail:     nil,
		leftNode: nil,
		tree:     tree,
	}
	return lru, nil
}

func (lru *LRU) Set(key string, obj interface{}) {
	lru.l.Lock()
	defer lru.l.Unlock()

	// if exists
	node := lru.get(key)
	if node != nil {
		node.value = obj
		return
	}

	// has left
	if lru.left > 0 {
		lru.left--

		if lru.leftNode == nil {
			node = new(Node)
		} else {
			node = lru.leftNode
			lru.leftNode = node.next
			node.clear()
		}

	} else {
		// not left, should del
		node = lru.tail
		lru.pick(node)
		lru.tree.Del(node.key)
		node.clear()
	}

	node.init(key, obj)
	lru.tree.Set(key, node)
	lru.push(node)

	return
}

func (lru *LRU) Get(key string) (obj interface{}) {
	lru.l.Lock()
	defer lru.l.Unlock()
	node := lru.get(key)
	if node == nil {
		return nil
	}
	return node.value
}

func (lru *LRU) get(key string) (node *Node) {
	node = lru.tree.Get(key)
	if node == nil {
		return nil
	}
	lru.pick(node)
	lru.push(node)
	return
}

func (lru *LRU) pick(node *Node) {
	if node == lru.head {
		lru.head = node.next
		lru.head.pre = nil
	} else {
		node.pre.next = node.next
		if node.next != nil {
			node.next.pre = node.pre
		}
	}
	if node == lru.tail {
		lru.tail = node.pre
		if lru.tail != nil {
			lru.tail.next = nil
		}
	} else {
		node.next.pre = node.pre
		if node.pre != nil {
			node.pre.next = node.next
		}
	}
	node.pre = nil
	node.next = nil
	return
}

func (lru *LRU) push(node *Node) {
	node.next = lru.head
	if lru.head != nil {
		lru.head.pre = node
	}
	lru.head = node
	if lru.tail == nil {
		lru.tail = node
	}
	return
}

func (lru *LRU) Del(key string) (obj interface{}) {
	lru.l.Lock()
	defer lru.l.Unlock()

	node := lru.tree.Del(key)
	if node == nil {
		return nil
	}

	lru.pick(node)

	node.clear()
	lru.left++
	node.next = lru.leftNode
	lru.leftNode = node

	return
}

func (lru *LRU) Print() {
	node := lru.head
	for node != nil {
		fmt.Printf("%+v\t", node.value)
		node = node.next
	}
	fmt.Println()
}
