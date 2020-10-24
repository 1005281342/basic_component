package main

import (
	"fmt"
	"github.com/1005281342/basic_component/cache"
	"log"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	case1()
	case2()
}

func case2() {
	var ca = cache.NewSimpleCache(cache.Opt{Callback: nil, DefaultExpiration: time.Second})
	ca.PutWithExpire(5, struct{}{}, 100*time.Millisecond)
	ca.PutWithExpire(6, struct{}{}, 100*time.Millisecond)
	ca.PutWithExpire(7, struct{}{}, 1000*time.Millisecond)
	var ok bool
	_, ok = ca.Get(5)
	fmt.Println(ok)
	time.Sleep(100 * time.Millisecond)
	_, ok = ca.Get(6)
	fmt.Println(ok)
	_, ok = ca.Get(7)
	fmt.Println(ok)
	ca.Remove(7)
	_, ok = ca.Get(7)
	fmt.Println(ok)
}

func case1() {
	var size = 1024 * 32
	var ca = cache.NewSimpleCache(cache.Opt{})

	var ch = make(chan struct{}, 1000)
	for i := 0; i < size*2+100; i++ {
		ch <- struct{}{}
		wg.Add(1)
		go func(a int) {
			defer func() {
				<-ch
				wg.Done()
			}()
			ca.Put(a, struct{}{})
		}(i)
	}
	wg.Wait()
	log.Println(ca.Len(), size)

	var cnt int
	for i := 0; i < size*2+100; i++ {
		if i&1 == 0 {
			ca.Remove(i)
		}

		if _, ok := ca.Get(i); ok {
			cnt++
			//log.Printf("err elem: %d", i)
		}
	}
	log.Printf("cnt: %d", cnt)
}
