package main

import "fmt"

func main() {
	demo1()
	demo2()
}

type a struct {
	b int
	c string
}

type aa struct {
	a
	d int
}

type aaa struct {
	*a
	d int
}

func demo1() {
	var ax = &aa{}
	ax.c = "xx"
	edit(ax)
	fmt.Println(ax.c)
}

func edit(aa2 *aa) {
	aa2.c = "xx2"
}

func demo2() {
	var ax = &aaa{a: &a{}}
	ax.c = "xx"
	edit2(ax)
	fmt.Println(ax.c)
}
func edit2(aa2 *aaa) {
	aa2.c = "xx2"
}
