package main

import (
	"fmt"

	"github.com/probably/filters"
)

func main() {
	bloomFilter := filters.NewBloomFilter(4, 0.001)
	bloomFilter.Add([]byte("hello"))
	bloomFilter.Add([]byte("world"))
	fmt.Println(bloomFilter.Contains([]byte("another world"))) // false
}
