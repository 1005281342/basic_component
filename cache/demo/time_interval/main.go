package main

import (
	"log"
	"time"
)

func main() {
	var (
		now   = time.Now()
		now0  = now.UnixNano()
		now5  = now.Add(time.Second * 5).UnixNano()
		now10 = now.Add(time.Second * 10).UnixNano()
		now12 = now.Add(time.Second * 12).UnixNano()
	)

	log.Println(x(now0, now0))
	log.Println(x(now0, now5))
	log.Println(x(now0, now10))
	log.Println(x(now0, now12))
}

func x(a, b int64) int {
	var in = time.Second * 10
	return int(float64(b-a) / float64(in))
}
