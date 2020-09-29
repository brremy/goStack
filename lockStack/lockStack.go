package lockStack

// ---------------------------------------------------------
// Reader/Writer lock stack of integers. For simplicity
// -1 is returned when peeking/popping from stack.
//
// The rwmutex is a clssic r/w lock that allows concurrent
// reads allows onle single writes. Incomming reads are
// block if there is a waiting writer.
//
// ---------------------------------------------------------

import (
	"sync"
)

type Node struct {
	next *Node
	val  int
}

type LockStack struct {
	head  *Node
	mutex sync.RWMutex
}

// Push an integer onto the stack.
//
func (curStack *LockStack) Push(value int) {
	curStack.mutex.Lock()
	defer curStack.mutex.Unlock()
	curStack.head = &Node{curStack.head, value}
}

// Pop an integer from the stack.
// Returns -1 if the stack is empty.
//
func (curStack *LockStack) Pop() int {
	curStack.mutex.Lock()
	defer curStack.mutex.Unlock()
	value := -1
	if curStack.head != nil {
		value = curStack.head.val
		curStack.head = curStack.head.next
	}

	return value
}

// Peek the top integer from the stack.
// Returns -1 if the stack is empty.
//
func (curStack *LockStack) Peek() int {
	curStack.mutex.RLock()
	defer curStack.mutex.RUnlock()
	value := -1
	if curStack.head != nil {
		value = curStack.head.val
	}

	return value
}
