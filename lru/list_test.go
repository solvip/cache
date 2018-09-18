package lru

import (
	"testing"
)

func TestMoveToHead(t *testing.T) {
	var l list

	n1 := &node{key: "n1", value: 1}
	n2 := &node{key: "n2", value: 2}
	n3 := &node{key: "n3", value: 3}

	l.moveToHead(n1)
	if l.head != n1 && l.tail != n1 {
		t.Fatalf("expected head and tail to point to n1; instead head: %#v, tail: %#v", l.head, l.tail)
	}
	if n1.next != nil && n1.prev != nil {
		t.Fatalf("expected n1.next = nil, n1.prev = nil; instead n1.next = %#v, n1.prev = %#v",
			n1.next, n1.prev)
	}

	l.moveToHead(n2)
	if l.head != n2 && l.tail != n1 {
		t.Fatalf("epxected l.head = &n2; l.tail = &n1; instead l.head = %#v, l.tail = %#v",
			l.head, l.tail)
	}
	if n2.prev != nil && n2.next != n1 {
		t.Fatalf("expected n2.prev = nil; n2.next = &n1; instead n2.prev = %#v, n2.next = %#v",
			n2.prev, n2.next)
	}
	if n1.prev != n2 && n1.next != nil {
		t.Fatalf("expected n1.prev = &n2, n1.next = nil; instead n1.prev = %#v, n1.next = %#v",
			n1.prev, n1.next)
	}

	// Put 3 to head; then move n1 to front.
	// The list should then be n1 -> n3 -> n2
	l.moveToHead(n3)
	l.moveToHead(n1)

	if l.head != n1 {
		t.Fatalf("expected l.head = n1; instead got %#v", l.head)
	}

	if l.tail != n2 {
		t.Fatalf("expected l.tail = n2; instead got %#v", l.tail)
	}

	if l.head != n1 && l.head.next != n3 && l.head.next.next != n2 && l.head.next.next.next != nil {
		t.Fatalf("expected next sequence from head to be n1 -> n3 -> n2 -> nil")
	}

	if l.tail != n2 && l.tail.prev != n1 && l.tail.prev.prev != n3 && l.tail.prev.prev.prev != nil {
		t.Fatalf("expected prev sequence from tail to be n2 -> n3 -> n1 -> nil")
	}
}

func TestDropTail(t *testing.T) {
	var l list

	n1 := &node{key: "n1", value: 1}
	n2 := &node{key: "n2", value: 2}
	n3 := &node{key: "n3", value: 3}

	if n := l.dropTail(); n != nil {
		t.Fatalf("expected dropTail() to return nil when list is empty")
	}

	l.moveToHead(n1)
	if n := l.dropTail(); n != n1 {
		t.Fatalf("expected dropTail()=n1; instead got %#v", n1)
	} else if n.next != nil || n.prev != nil {
		t.Fatalf("expected dropTail() to set next/prev of returned node to nil")
	} else if l.head != nil || l.tail != nil {
		t.Fatalf("expected dropTail() to set head/tail to nil when last item is dropped")
	}

	l.moveToHead(n1)
	l.moveToHead(n2)
	l.moveToHead(n3)

	if n := l.dropTail(); n != n1 {
		t.Fatalf("expected dropTail()=n1; instead got %#v", n1)
	} else if l.tail != n2 {
		t.Fatalf("expected dropTail() to set l.tail = n2; instead l.tail = %#v", l.tail)
	} else if l.head != n3 {
		t.Fatalf("expected dropTail() to leave l.head = n3; instead l.head = %#v", l.head)
	}
}
