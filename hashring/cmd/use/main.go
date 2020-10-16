package main

import (
	"fmt"
	"github.com/1005281342/basic_component/hashring"
)

func main() {
	var hr = hashring.New(10, nil)
	_ = hr.Add("kkkx")
	_ = hr.Add("aasfda")
	_ = hr.Add("fafcxz")
	var x string
	x, _ = hr.Locate("kkkx:9")
	fmt.Println(x)
	x, _ = hr.Locate("fffxaz")
	fmt.Println(x)
	x, _ = hr.Locate("ggg0")
	fmt.Println(x)
	x, _ = hr.Locate("ggg1")
	fmt.Println(x)
	x, _ = hr.Locate("ggg2")
	fmt.Println(x)
	x, _ = hr.Locate("ggg3")
	fmt.Println(x)
	x, _ = hr.Locate("ggg4")
	fmt.Println(x)
	x, _ = hr.Locate("ggg5")
	fmt.Println(x)
	x, _ = hr.Locate("ggg6")
	fmt.Println(x)
	x, _ = hr.Locate("ggg7")
	fmt.Println(x)
	x, _ = hr.Locate("ggg8")
	fmt.Println(x)
	x, _ = hr.Locate("ggg9")
	fmt.Println(x)
}
