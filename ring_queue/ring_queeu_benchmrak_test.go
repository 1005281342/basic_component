package ring_queue

import "testing"

func BenchmarkRingQueue_Insert(b *testing.B) {
	var rq = NewRingQueue(1024)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rq.Insert(i)
	}
}

func BenchmarkRingQueue_LInsert(b *testing.B) {
	var rq = NewRingQueue(1024)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rq.LInsert(i)
	}
}

func BenchmarkRingQueue_Pop(b *testing.B) {
	var rq = NewRingQueue(1024)
	for i := 0; i < 1024; i++ {
		rq.Insert(i)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rq.Pop()
	}
}
