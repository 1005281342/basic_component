package old_lfu

import (
	"reflect"
	"sync"
	"testing"
)

func TestCache_Get(t *testing.T) {

	type fields struct {
		cache    *sync.Map
		freqMap  *sync.Map
		size     uint32
		capacity uint32
		min      uint32
	}
	type args struct {
		key interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		{"get !ok", fields{}, args{key: 10}, nil},
		{"get ok", fields{}, args{key: uint32(1)}, uint32(3)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewLFUCache(3)
			c.Put(uint32(1), uint32(2))
			c.Put(uint32(1), uint32(3))
			if got, _ := c.Get(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_GetMin(t *testing.T) {
	type fields struct {
		cache    *sync.Map
		freqMap  *sync.Map
		size     uint32
		capacity uint32
		min      uint32
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{"ok", fields{}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewLFUCache(3)
			c.Put(12, 2)
			c.Put(12, 2)
			c.Put(12, 2)
			c.Put(123, 2)
			c.Put(123, 2)
			c.Put(123, 2)
			c.Put(12223, 2)
			c.Put(12223, 2)
			c.Put(12223, 2)
			c.Put(1213, 2)
			//c.Put(1234, 2)
			//c.Put(1234, 2)
			//c.Put(1222223, 2)
			if got := c.GetMin(); got != tt.want {
				t.Errorf("GetMin() = %v, want %v", got, tt.want)
			}
			t.Log(c.Get(123))
			t.Log(c.Get(12223))
			t.Log(c.Get(12))
		})
	}
}

func TestCache_Put(t *testing.T) {
	type fields struct {
		cache    *sync.Map
		freqMap  *sync.Map
		size     uint32
		capacity uint32
		min      uint32
	}
	type args struct {
		key   interface{}
		value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"capacity == 0", fields{}, args{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewLFUCache(0)
			c.Put(1, 2)
			if x, _ := c.Get(1); x != nil {
				t.Fail()
			}
		})
	}
}
