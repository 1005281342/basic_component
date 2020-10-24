package main

import (
	"github.com/1005281342/basic_component/cache"
	"log"
	"sync"
)

var wg sync.WaitGroup

func main() {
	var size = 1024 * 32
	var lfu = cache.NewLFUCache(&cache.Opt{Capacity: size})

	var ch = make(chan struct{}, 1000)
	for i := 0; i < size*2+100; i++ {
		ch <- struct{}{}
		wg.Add(1)
		go func(a int) {
			defer func() {
				<-ch
				wg.Done()
			}()
			lfu.Put(a, struct{}{})
		}(i)
	}
	wg.Wait()
	log.Println(lfu.Len(), size)

	var cnt int
	for i := 0; i < size*2+100; i++ {
		if _, ok := lfu.Get(i); ok {
			cnt++
			//log.Printf("err elem: %d", i)
		}
	}
	log.Printf("cnt: %d", cnt)
}
