package main

// ---------------------------------------------------------
// unit test and benchmark driver
//
// ---------------------------------------------------------
import (
	"fmt"
	"flag"
	"math/rand"
	"time"
	
	"../occStack"
	"../lockStack"
)

// Interface for abstracting the tests.
//
type ConcurrentStack interface {
	Push(value int)
	Pop() int
	Peek() int
}

// Test peek on an empty stack
//
func testPeekEmpty(stack ConcurrentStack) {
	value := stack.Peek()
	if (value != -1) {
		panic(fmt.Sprintf("Expected value 1, received %d", value))
	}
}

// Test pop on an empty stack
//
func testPopEmpty(stack ConcurrentStack) {
	value := stack.Pop()
	if (value != -1) {
		panic(fmt.Sprintf("Expected value 1, received %d", value))
	}
}

// Test pop on an empty stack
//
func testBasicPushPop(stack ConcurrentStack) {
	stack.Push(1)
	value := stack.Pop()
	if (value != 1) {
		panic(fmt.Sprintf("Expected value 1, received %d", value))
	}
}

// Test pop on an empty stack
//
func testInterleavedPushPop(stack ConcurrentStack) {
	stack.Push(1)
	stack.Push(2)
	value := stack.Pop()
	if (value != 2) {
		panic(fmt.Sprintf("Expected value 2, received %d", value))
	}

	stack.Push(3)
	value = stack.Pop()
	if (value != 3) {
		panic(fmt.Sprintf("Expected value 3, received %d", value))
	}

	stack.Push(4)
	value = stack.Pop()
	if (value != 4) {
		panic(fmt.Sprintf("Expected value 4, received %d", value))
	}

	value = stack.Pop()
	if (value != 1) {
		panic(fmt.Sprintf("Expected value 1, received %d", value))
	}
}

// Test some combined functionality.
//
func testCombined(stack ConcurrentStack) {
	stack.Push(5)
    stack.Push(4)
	stack.Push(3)
	value := stack.Peek()
	if (value != 3) {
		panic(fmt.Sprintf("Expected value 3, received %d", value))
	}

	value = stack.Pop()
	if (value != 3) {
		panic(fmt.Sprintf("Expected value 3, received %d", value))
	}

	value = stack.Pop()
	if (value != 4) {
		panic(fmt.Sprintf("Expected value 4, received %d", value))
	}

	value = stack.Pop()
	if (value != 5) {
		panic(fmt.Sprintf("Expected value 5, received %d", value))
	}

	value = stack.Pop()
	if (value != -1) {
		panic(fmt.Sprintf("Expected value -1, received %d", value))
	}
}

// Benchmark the stacks using the given degrees of parallelism, write percent, and test duration.
// The benchmark spawns a number of threads specified by dop and each thread randomly performs
// read and write operations specified by writePercent until the test terminates
// 
func benchmark(stack ConcurrentStack, dop int, writePercent int, testDuration int) int64 {
	terminateChan := make(chan int)
	resultChan := make(chan int64)
	for i := 0; i <= dop; i++ {
		go func() {
			// track the troughput
			//
			var curThroughput int64 = 0

			// track the number of local trhead push/pop to make sure we only pop when we have a push
			//
			var addedStackSize int64 = 0

			for {

				// Send the throughput for aggregation and terminate the go rutine
				//
				select {
				case _, chanStillOpen := <-terminateChan:
					if (!chanStillOpen) {
						resultChan <- curThroughput
						return
					} else {
						panic("Unexpected message on terminate channel")
					}

				default: // Nothing here, continue.
				}

				// generate a random number from [0 - 99]
				// Note calling rand every time and not skipping when 
				// writePercent == 0 or writePercent == 1 keeps the performance
				// consistent
				//
				rwRand := rand.Intn(99)

				if (rwRand >= writePercent-1) {
					// perform a read operation
					//
					stack.Peek()
				} else {
					// perfrom a write operation. We want to call this every
					//
					pushPopRand := rand.Intn(1)
					if (pushPopRand == 0 && addedStackSize > 0) {
						stack.Pop()
						addedStackSize--
					} else {
						stack.Push(1)
						addedStackSize++
					}
				}

				curThroughput++
			}
		}()
	}

	// Sleep for the test duration and terminate the worker threads.
	//
	go func() {
		time.Sleep(time.Duration(testDuration) * time.Second)
		
		close(terminateChan)
    }()

	var throughput int64 = 0
	for i := 0; i <= dop; i++ {
		throughput += <-resultChan
	}
	
	return throughput
}

func main() {
	testPtr := flag.Bool("unitTest", false, "run unit tests")
	dopPtr := flag.Int("dop", 1, "integer degrees of parallism.")
	writePercentPtr := flag.Int("writePercent", 0, "integer percent of write querries out of 100.")
	durationPtr := flag.Int("duration", 5, "integer benchmark duration seconds.")
	fullBenchmarkPtr := flag.Bool("fullBenchmark", false, "run full benchmark")
	flag.Parse()

	if (*testPtr) {
		for i := 0; i < 2; i++ {
			var stack ConcurrentStack
			if (i==0) {
				stack  = new(occStack.OccStack)
			} else {
				stack  = new(lockStack.LockStack)
			}

			testBasicPushPop(stack)
			testPeekEmpty(stack)
			testPopEmpty(stack)
			testInterleavedPushPop(stack)
			testCombined(stack)
		}

		fmt.Print("Unit tests succeeded.\n")
	} else if(!*fullBenchmarkPtr) {
		stack := new(occStack.OccStack)
		occThroughput := benchmark(stack, *dopPtr, *writePercentPtr, *durationPtr)
		fmt.Printf("OCC troughput: %d\n", occThroughput)

		stack = new(lockStack.LockStack)
		lockThroughput := benchmark(stack, *dopPtr, *writePercentPtr, *durationPtr)
		fmt.Printf("Locking troughput: %d\n", lockThroughput)
	}
	else{
		for dop :=1; dop <=10; dop++
		{
			for writePercent := 0; writePercent <= 100; writePercent += 10 {
				stack := new(occStack.OccStack)
				occThroughput := benchmark(stack, *dopPtr, *writePercentPtr, *durationPtr)
				fmt.Printf("OCC, dop: %d, writePercent: %d, troughput: %d\n", dop, writePercent, occThroughput)

				stack = new(lockStack.LockStack)
				lockThroughput := benchmark(stack, *dopPtr, *writePercentPtr, *durationPtr)
				fmt.Printf("Locking, dop: %d, writePercent: %d, troughput: %d\n", dop, writePercent, lockThroughput)
			}
			
			fmt.Printf("\n")
		}
	}
}
