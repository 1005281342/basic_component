package main

import (
	"fmt"
	"sync"
)

func main() {
	var mp sync.Map
	for i := 0; i < 16; i++ {
		mp.Store(i, struct{}{})
	}

	mp.Range(func(key, value interface{}) bool {
		fmt.Println(key)
		mp.Delete(key)
		return true
	})

	mp.Range(func(key, value interface{}) bool {
		fmt.Println(key)
		//mp.Delete(key)
		return true
	})
}
