package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, filters!!!")
}

type BloomFilter struct {
	m uint64 // ?
}

func NewBloomFilter(m uint64) BloomFilter {
	return BloomFilter{m}
}

func (filter BloomFilter) Add(value string) {
	// TODO
}

func (filter BloomFilter) Contains(value string) bool {
	return false
}
