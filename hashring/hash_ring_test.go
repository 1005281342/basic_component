package hashring

import (
	"bufio"
	"fmt"
	"github.com/spaolacci/murmur3"
	"github.com/vedhavyas/hashring"
	"log"
	"os"
	"testing"
)

var (
	replicaCount = 3
	hashFunc     = murmur3.New32()
	hr           = hashring.New(replicaCount, hashFunc)
	thr          = hashring.New(replicaCount, hashFunc)
)

func BenchmarkHashRing_Add(b *testing.B) {

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := hr.Add(fmt.Sprintf("node-%d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func BenchmarkHashRing_T_Add(b *testing.B) {

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := thr.Add(fmt.Sprintf("node-%d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func BenchmarkHashRing_Locate(b *testing.B) {
	fd, _ := os.Open("/usr/share/dict/words")
	defer fd.Close()
	scanner := bufio.NewScanner(fd)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ok := scanner.Scan()
		if !ok {
			break
		}
		var text = scanner.Text()
		_, err := hr.Locate(text)
		if err != nil {
			log.Fatal(err)
		}
	}

	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}
}

func BenchmarkHashRing_T_Locate(b *testing.B) {
	fd, _ := os.Open("/usr/share/dict/words")
	defer fd.Close()
	scanner := bufio.NewScanner(fd)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ok := scanner.Scan()
		if !ok {
			break
		}
		var text = scanner.Text()
		_, err := thr.Locate(text)
		if err != nil {
			log.Fatal(err)
		}
	}

	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}
}

func BenchmarkHashRing_Delete(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := hr.Delete(fmt.Sprintf("node-%d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func BenchmarkHashRing_T_Delete(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := thr.Delete(fmt.Sprintf("node-%d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
}
