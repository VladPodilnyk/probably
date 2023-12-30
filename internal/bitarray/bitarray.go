package bitarray

import (
	"math"
)

type BitArray struct {
	data []byte
}

func New(size uint) *BitArray {
	sizeInBytes := uint(math.Ceil(float64(size) / 8))
	array := make([]byte, sizeInBytes)
	return &BitArray{data: array}
}

func (array *BitArray) Set(bit uint) {
	index := bit / 8
	array.data[index] |= 1 << (bit % 8)
}

func (array *BitArray) Merge(value *BitArray) {
	for i := range array.data {
		array.data[i] |= value.data[i]
	}
}

func (array *BitArray) Union(value *BitArray) *BitArray {
	result := make([]byte, len(array.data))
	for i := range array.data {
		result[i] = array.data[i] | value.data[i]
	}
	return &BitArray{data: result}
}

func (array *BitArray) IsSet(bit uint) bool {
	index := bit / 8
	return array.data[index]&(1<<(bit%8)) != 0
}

func (array *BitArray) Clear() {
	for i := range array.data {
		array.data[i] = 0
	}
}
