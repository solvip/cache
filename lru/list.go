package lru

type list struct {
	head *node
	tail *node
}

type node struct {
	key   string
	value int

	next *node
	prev *node
}

// moveToHead - move node to the head of the list
// If the node isn't a member of the list, then it's simply pushed to the head.
func (l *list) moveToHead(node *node) {
	if l.head == node {
		// If the node is already at the head of the list, we do nothing
		return
	} else if l.head == nil {
		// The node we're pushing now is to be the first node of the list.
		// It must become both the head and the tail of the list.
		l.head = node
		l.tail = node
		node.next = nil
		node.prev = nil

		return
	}

	// If the node is currently the tail; set the tail to the previous node.
	// This is safe as we've already returned early in the case of the head == node.
	if l.tail == node {
		l.tail = node.prev
	}

	if node.prev != nil {
		node.prev.next = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	}

	node.next = l.head
	l.head.prev = node
	l.head = node

	return
}

// dropTail - remove the tail node from the list
func (l *list) dropTail() (prevTail *node) {
	if l.tail == nil {
		return nil
	}

	prevTail = l.tail
	l.tail = l.tail.prev
	prevTail.next = nil
	prevTail.prev = nil
	if l.tail == nil {
		// We just dropped the last item.
		// Unset the head pointer as well.
		l.head = nil
		return prevTail
	}
	l.tail.next = nil

	return prevTail
}
