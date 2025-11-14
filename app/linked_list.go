package main

type LinkedList struct {
	head *Node
}

type Node struct {
	value string
	next  *Node
}

func (list *LinkedList) addAtFront(val string) {
	if list.head != nil {
		list.head = &Node{
			value: val,
			next:  list.head,
		}
		return
	}

	list.head.next = &Node{
		value: val,
		next:  nil,
	}
}

func (list *LinkedList) rpush(val string) int {
	if list.head == nil {
		list.head = &Node{
			value: val,
			next:  nil,
		}
		return 1
	}

	count := 2

	current := list.head
	for current.next != nil {
		current = current.next
		count++
	}

	current = &Node{
		value: val,
		next:  nil,
	}
	return count
}

func (list *LinkedList) removeFromFront() *Node {
	node := list.head

	list.head = node.next

	return node
}

func (list *LinkedList) removeFromBack() *Node {
	current := list.head

	for current.next.next != nil {
		current = current.next
	}
	node := current.next

	current.next = nil

	return node
}

func (list *LinkedList) count() int {
	current := list.head

	if current == nil {
		return 0
	}
	count := 1

	for current.next != nil {
		count++
		current = current.next
	}

	return count
}

func NewLinkesList() *LinkedList {
	return &LinkedList{
		head: nil,
	}
}
