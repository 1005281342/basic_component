package ring_queue

import (
	"fmt"
	"testing"
)

func TestRingQueueBlock_Len(t *testing.T) {
	var r = NewRingQueueBlock(10)
	r.Insert(123)
	r.Insert(123)
	r.Insert(123)
	t.Logf("elem count: %d", r.Len()) // 3

	for i := 0; i < 7; i++ {
		r.Insert(i) // 3 + 7 -> 10
	}
	r.Pop() // 9
	fmt.Println(r.Head(), r.Tail())
	r.Pop()
	r.Pop()
	r.Pop()
	fmt.Println(r.Head(), r.Tail()) // 9-3 = 6

	r.Insert(1023)
	r.Insert(1024)
	fmt.Println(r.Head(), r.Tail()) // 8

	t.Logf("elem count: %d", r.Len()) // ok 8
}
