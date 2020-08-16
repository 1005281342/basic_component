package ring_queue_nonblocking

import (
	"reflect"
	"testing"
)

func TestNewRingQueue(t *testing.T) {

	var (
		CASE1 = "cap <= 1"
		OK    = "OK"
	)

	type args struct {
		cap int
	}
	tests := []struct {
		name string
		args args
		want *RingQueue
	}{
		{CASE1, args{cap: 1}, nil},
		{OK, args{cap: 2}, &RingQueue{cap: 2, queue: make([]interface{}, 2, 2)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRingQueue(tt.args.cap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRingQueue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRingQueue_Head(t *testing.T) {
	type fields struct {
		cap   int
		queue []interface{}
		index int
		head  int
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{"", fields{
			cap:   5,
			index: 4,
			head:  0,
			queue: []interface{}{1, 2, 3, 4, 5},
		}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RingQueue{
				cap:   tt.fields.cap,
				queue: tt.fields.queue,
				index: tt.fields.index,
				head:  tt.fields.head,
			}
			if got := r.Head(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Head() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRingQueue_Insert(t *testing.T) {
	type fields struct {
		cap   int
		queue []interface{}
		index int
		head  int
	}
	type args struct {
		x interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		{"", fields{
			5, []interface{}{1, 2, 3, 4, 5}, 4, 0,
		}, args{6}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RingQueue{
				cap:   tt.fields.cap,
				queue: tt.fields.queue,
				index: tt.fields.index,
				head:  tt.fields.head,
			}
			if got := r.Insert(tt.args.x); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Insert() = %v, want %v", got, tt.want)
			}
			t.Logf("now head: %v, tail: %v", r.Head(), r.Tail())
		})
	}
}

func TestRingQueue_IsFull(t *testing.T) {
	type fields struct {
		cap   int
		queue []interface{}
		index int
		head  int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"", fields{
			5, []interface{}{1, 2, 3, 4, 5}, 4, 0,
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RingQueue{
				cap:   tt.fields.cap,
				queue: tt.fields.queue,
				index: tt.fields.index,
				head:  tt.fields.head,
			}
			if got := r.IsFull(); got != tt.want {
				t.Errorf("IsFull() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRingQueue_LInsert(t *testing.T) {
	type fields struct {
		cap   int
		queue []interface{}
		index int
		head  int
	}
	type args struct {
		x interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		{"full ", fields{
			5, []interface{}{1, 2, 3, 4, 5}, 4, 0,
		}, args{6}, 5},
		{"not full", fields{
			5, []interface{}{1, 2, 3, nil, nil}, 2, 0,
		}, args{6}, nil},
		{"index == 0 and full", fields{
			5, []interface{}{1, 2, 3, 4, 5}, 0, 1,
		}, args{6}, 1},
		{"a: index != 0 and not full", fields{
			5, []interface{}{1, nil, 3, 4, 5}, 0, 2,
		}, args{6}, nil},
		{"b: index == 0 and not full", fields{
			5, []interface{}{1, 2, 3, 4, nil}, 3, 0,
		}, args{6}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RingQueue{
				cap:   tt.fields.cap,
				queue: tt.fields.queue,
				index: tt.fields.index,
				head:  tt.fields.head,
			}
			if got := r.LInsert(tt.args.x); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LInsert() = %v, want %v", got, tt.want)
			}
			t.Logf("now head: %v, tail: %v", r.Head(), r.Tail())
		})
	}
}

func TestRingQueue_LPop(t *testing.T) {
	type fields struct {
		cap   int
		queue []interface{}
		index int
		head  int
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{"", fields{
			5, []interface{}{nil, nil, nil, nil, nil}, 0, 0,
		}, nil},
		{"", fields{
			5, []interface{}{1, 2, 3, 4, 5}, 4, 0,
		}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RingQueue{
				cap:   tt.fields.cap,
				queue: tt.fields.queue,
				index: tt.fields.index,
				head:  tt.fields.head,
			}
			if got := r.LPop(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LPop() = %v, want %v", got, tt.want)
			}
			t.Logf("now head %v, tail %v, len %v", r.Head(), r.Tail(), r.Len())
		})
	}
}

func TestRingQueue_Pop(t *testing.T) {
	type fields struct {
		cap   int
		queue []interface{}
		index int
		head  int
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{"", fields{
			5, []interface{}{nil, nil, nil, nil, nil}, 0, 0,
		}, nil},
		{"", fields{
			5, []interface{}{1, 2, 3, 4, 5}, 4, 0,
		}, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RingQueue{
				cap:   tt.fields.cap,
				queue: tt.fields.queue,
				index: tt.fields.index,
				head:  tt.fields.head,
			}
			if got := r.Pop(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pop() = %v, want %v", got, tt.want)
			}
			t.Logf("now head %v, tail %v, len %v", r.Head(), r.Tail(), r.Len())
		})
	}
}

func TestRingQueue_Tail(t *testing.T) {
	type fields struct {
		cap   int
		queue []interface{}
		index int
		head  int
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{"", fields{
			cap: 5, head: 0, index: 4,
			queue: []interface{}{1, 2, 3, 4, 5},
		}, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RingQueue{
				cap:   tt.fields.cap,
				queue: tt.fields.queue,
				index: tt.fields.index,
				head:  tt.fields.head,
			}
			if got := r.Tail(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRingQueue_Len(t *testing.T) {
	type fields struct {
		cap   int
		queue []interface{}
		index int
		head  int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{"head <= tail", fields{
			cap: 4, index: 1, head: 0, queue: []interface{}{1, 2, nil, nil},
		}, 2},
		{"head > tail", fields{
			cap: 5, index: 0, head: 3, queue: []interface{}{1, nil, nil, 4, 5},
		}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RingQueue{
				cap:   tt.fields.cap,
				queue: tt.fields.queue,
				index: tt.fields.index,
				head:  tt.fields.head,
			}
			if got := r.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}
