package main

import "fmt"

type entryWithFreq struct {
	entry
	freq int
}

func (e *entryWithFreq) Reset() {
	e.entry.Reset()
	e.freq = 0
}

type entry struct {
	key interface{}
	item
}

type item struct {
	value      interface{} // 元素值
	expiration int64       // 绝对过期时间
}

func (i *item) Reset() {
	i.value = nil
	i.expiration = 0
}

func (e *entry) Reset() {
	e.item.Reset()
	e.key = nil
}

func main() {
	var e = entryWithFreq{freq: 100,
		entry: entry{key: 100, item: item{value: 100, expiration: 100}}}
	fmt.Printf("%+v \n", e)
	e.Reset()
	fmt.Printf("%+v \n", e)
}
