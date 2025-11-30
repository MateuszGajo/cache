package main

type LinkedList struct {
	head *Node
	tail *Node
	size int
}

type Node struct {
	value string
	next  *Node
	prev  *Node
}

func NewNode(val string) *Node {
	return &Node{
		value: val,
	}
}

// add lpush
// it will always be header
// curent head should be moved to .next
// currnet head should add prev to new head node
// increment size by one

func (list *LinkedList) lpush(val string) int {
	node := NewNode(val)
	currentHead := list.head

	list.head = node
	if currentHead != nil {
		currentHead.prev = node
		list.head.next = currentHead
	}

	list.size++

	return list.size
}

func (list *LinkedList) Lpop(count int) []string {
	data := []string{}

	for i := 0; i < count; i++ {
		data = append(data, list.lpop())
	}

	return data
}

func (list *LinkedList) lpop() string {
	currentHead := list.head

	if currentHead == nil {
		return ""
	}

	next := currentHead.next

	list.head = next
	if next != nil {
		next.prev = nil
	}
	list.size--

	return currentHead.value
}

// blopop
// command wait for a time specifried or undeffined amount if list is empty
// so check if list is not empty reutnr element if empty do a ticketing with like checking every 50ms?
// next phase addiitonal timeout, in timeout condition also check if elements exists if not then return null

func (list *LinkedList) rpush(val string) int {
	node := NewNode(val)
	if list.head == nil {
		list.head = node
		list.tail = node
		list.size = 1
		return 1
	}

	count := 2

	current := list.head
	for current.next != nil {
		current = current.next
		count++
	}

	node.prev = current

	current.next = node
	list.tail = node
	list.size = count
	return count
}

func (list *LinkedList) getRange(start int, end int) []string {
	if start < 0 && end < 0 {
		// more efficient lookup starting from tail
		return list.getRangeNegativeLookup(start, end)
	}

	if start < 0 {
		start = list.size + start
	}

	if end < 0 {
		end = list.size + end
	}

	if start > end {
		return []string{}
	}

	return list.getRangePositiveLookup(start, end)
}

func (list *LinkedList) getRangePositiveLookup(start int, end int) []string {
	current := list.head
	i := 0
	output := []string{}
	for current != nil && i < start {
		current = current.next
		i++
	}

	for current != nil && i <= end {
		output = append(output, current.value)
		current = current.next
		i++
	}

	return output
}

// input are e.g from -2 to -1
func (list *LinkedList) getRangeNegativeLookup(start int, end int) []string {
	// we need to inverted to then do previous lookup
	// -1 becasue first element is -1 not 0
	startIverted := (end * -1) - 1
	endInverted := (start * -1) - 1
	current := list.tail
	i := 0
	output := []string{}
	for current != nil && i < startIverted {
		current = current.prev
		i++
	}

	for current != nil && i <= endInverted {
		output = append(output, current.value)
		current = current.prev
		i++
	}

	reverseStrings(output)

	return output
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
