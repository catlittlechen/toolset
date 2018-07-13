// Author: catlittlechen@gmail.com

package uqueue

import (
	"fmt"
	"sync"
)

type UQueue struct {
	l    *sync.Mutex
	m    map[string]struct{}
	head *Node
	tail *Node
}

type Value interface {
	UniuqID() string
}

type Node struct {
	value Value
	next  *Node
}

func New() *UQueue {
	uq := &UQueue{
		l:    new(sync.Mutex),
		m:    make(map[string]struct{}),
		head: nil,
		tail: nil,
	}
	return uq
}

func (uq *UQueue) Push(value Value) {
	uq.l.Lock()
	defer uq.l.Unlock()
	key := value.UniuqID()
	if _, ok := uq.m[key]; ok {
		return
	}
	node := &Node{
		value: value,
		next:  nil,
	}

	if uq.tail == nil {
		uq.tail = node
		uq.head = uq.tail
	} else {
		uq.tail.next = node
		uq.tail = node
	}
	return
}

func (uq *UQueue) Pop() Value {
	uq.l.Lock()
	defer uq.l.Unlock()

	node := uq.head
	if node == nil {
		return nil
	}

	uq.head = uq.head.next
	if uq.head == nil {
		uq.tail = nil
	}

	key := node.value.UniuqID()
	delete(uq.m, key)
	return node.value
}

func (uq *UQueue) Print() {
	uq.l.Lock()
	defer uq.l.Unlock()

	node := uq.head
	for node != nil {
		fmt.Print(node.value.UniuqID() + "\t")
		node = node.next
	}
	fmt.Println()
}
