package occStack

// ---------------------------------------------------------
// Optimistic concurrent stack of integers. For simplicity
// -1 is returned when peeking/popping from stack
//
// ---------------------------------------------------------

import (
	"sync/atomic"
	"unsafe"
)

type Node struct {
    next *Node
    val  int
}

type OccStack struct {
    head *Node
}

// Push an integer onto the stack.
//
func (curStack *OccStack) Push(value int) {
	for {
		// synchronization point.
		//
		curHead := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&curStack.head))))
		newHead := &Node{curHead, value}

		// Note: we are relying on the golang garbage collector to not reuse
		// a memory address while it is still accessable from another thread
		// this avoids the ABA problem, https://en.wikipedia.org/wiki/ABA_problem.
		//
		result := atomic.CompareAndSwapPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&curStack.head)),
			unsafe.Pointer(curHead),
			unsafe.Pointer(newHead))
		if (!result) {
			// retry
			continue
		}

		break
	}
}

// Pop an integer from the stack.
// Returns -1 if the stack is empty.
//
func (curStack *OccStack) Pop() int {
	value := -1
	
	for {
		value = -1
		// synchronization point.
		//
		curHead := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&curStack.head))))

		if (curHead != nil) {
			// We are relying on the golang garbage collector to not deallocate this
			// memory on concurrent pops because it is still accessable from the
			// the current thread.
			//
			value = curHead.val
			result := atomic.CompareAndSwapPointer(
				(*unsafe.Pointer)(unsafe.Pointer(&curStack.head)),
				unsafe.Pointer(curHead),
				unsafe.Pointer(curHead.next))
			if (!result) {
				// retry
				continue
			}
		}

		break
	}

    return value
}

// Peek the top integer from the stack.
// Returns -1 if the stack is empty.
//
func (curStack *OccStack) Peek() int {
	value := -1
	
	// synchronization point.
	//
	curHead := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&curStack.head))))
	if (curHead != nil) {
        value = curHead.val
    }

    return value
}
