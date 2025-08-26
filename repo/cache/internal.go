package cache

import (
	"fmt"
	"strings"
)

type node struct {
	val  []byte
	next *node
	prev *node
}

// use with list.len--
func (n *node) Destructor() error {
	if n.next == nil || n.prev == nil {
		return fmt.Errorf("LRU error, cant remove tail or head")
	}
	n.prev.next = n.next
	n.next.prev = n.prev
	return nil
}

type list struct {
	head *node
	tail *node
	len  int
}

func NewList() list {
	head := node{nil, nil, nil}
	tail := node{nil, nil, &head}
	head.next = &tail
	return list{&head, &tail, 0}
}

func (l *list) AppendToHead(val []byte) *node {
	newNode := node{val, nil, nil}
	l.head.next.prev = &newNode
	newNode.next = l.head.next
	l.head.next = &newNode
	newNode.prev = l.head
	l.len++
	return &newNode
}

func (l *list) AppendToTail(val []byte) {
	newNode := node{val, nil, nil}
	l.tail.prev.next = &newNode
	newNode.prev = l.tail.prev
	l.tail.prev = &newNode
	newNode.next = l.tail
	l.len++
}

func (l *list) RemoveLast() {
	l.tail.prev.prev.next = l.tail
	l.tail.prev = l.tail.prev.prev
	l.len--
}

func (l *list) Represent() string {
	node := l.head.next
	res := []string{}
	for node.next != nil {
		res = append(res, string(node.val))
		node = node.next
	}
	return strings.Join(res, " | ")
}
