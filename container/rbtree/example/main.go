package main

import (
	"fmt"
	"github.com/1005281342/basic_component/container/rbtree"
)

func g(n int) {
	fmt.Printf("check: %v\n", rbtree.Check(n))
}

func main() {
	g(20)
	g(1024)
	g(8888)
}
