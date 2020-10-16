package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/spaolacci/murmur3"
	"github.com/vedhavyas/hashring"
)

var (
	replicaCount = flag.Int("rc", 4, "replication count")
	nodeCount    = flag.Int("nc", 8, "node count")
	keyCount     = flag.Int("kc", 6000000, "key count")
)

func main() {
	flag.Parse()

	var cnt int

	var hash = murmur3.New32()
	hr := hashring.New(*replicaCount, hash)
	thr := hashring.New(*replicaCount, hash)

	nodeMap := make(map[string]int)
	for i := 0; i < *nodeCount; i++ {
		err := hr.Add(fmt.Sprintf("node-%d", i))
		if err != nil {
			log.Fatal(err)
		}

		err = thr.Add(fmt.Sprintf("node-%d", i))
		if err != nil {
			log.Fatal(err)
		}
	}

	fd, _ := os.Open("/usr/share/dict/words")
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	for i := 0; i < *keyCount; i++ {
		ok := scanner.Scan()
		if !ok {
			break
		}
		var text = scanner.Text()
		n, err := hr.Locate(text)
		if err != nil {
			log.Fatal(err)
		}

		tn, err := thr.Locate(text)
		if err != nil {
			log.Fatal(err)
		}

		if n != tn {
			cnt++
		}

		nodeMap[n]++
	}

	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}

	for k, v := range nodeMap {
		fmt.Printf("%s: %d\n", k, v)
	}
	fmt.Printf("cnt: %d", cnt)
}
