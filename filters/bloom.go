package filters

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"hash"
	"math"

	"github.com/probably/internal/bitarray"
)

const HASH_CLAMP_VALUE = 4

type BloomFilter struct {
	bitarray       *bitarray.BitArray
	numberOfBits   uint
	numberOfHashes uint
	filterSize     uint
	fpRate         float64
}

func NewBloomFilter(size uint, fpRate float64) BloomFilter {
	array_bits := calculateFilterSize(size, fpRate)
	array := bitarray.New(array_bits)
	hashes := calculateNumberOfHashes(size, array_bits)
	return BloomFilter{
		bitarray:       array,
		numberOfBits:   array_bits,
		numberOfHashes: hashes,
		filterSize:     size,
		fpRate:         fpRate,
	}
}

func (filter BloomFilter) Add(data []byte) {
	calculateHash := hashGenerator(data, filter.numberOfBits)
	for i := uint(0); i < filter.numberOfHashes; i++ {
		index := calculateHash(i)
		filter.bitarray.Set(index)
	}
}

func (filter BloomFilter) Contains(data []byte) bool {
	calculateHash := hashGenerator(data, filter.numberOfBits)
	for i := uint(0); i < filter.numberOfHashes; i++ {
		index := calculateHash(i)
		if !filter.bitarray.IsSet(index) {
			return false
		}
	}
	return true
}

func (filter BloomFilter) Merge(value BloomFilter) error {
	if filter.filterSize != value.filterSize || filter.fpRate != value.fpRate {
		return errors.New("BloomFilters must have the same configuration (size and fpRate)")
	}
	filter.bitarray.Merge(value.bitarray)
	return nil
}

func (filter BloomFilter) Union(value BloomFilter) (BloomFilter, error) {
	if filter.filterSize != value.filterSize || filter.fpRate != value.fpRate {
		return BloomFilter{}, errors.New("BloomFilters must have the same configuration (size and fpRate)")
	}
	union := filter.bitarray.Union(value.bitarray)
	return BloomFilter{
		bitarray:       union,
		numberOfBits:   filter.numberOfBits,
		numberOfHashes: filter.numberOfHashes,
		filterSize:     filter.filterSize,
		fpRate:         filter.fpRate,
	}, nil
}

func (filter BloomFilter) Size() uint {
	return filter.filterSize
}

func (filter BloomFilter) Clear() {
	filter.bitarray.Clear()
}

func hashGenerator(value []byte, hashMapSize uint) func(uint) uint {
	index := getHashValue(value, md5.New())
	offset := getHashValue(value, sha1.New())

	return func(i uint) uint {
		index = (index + offset) % hashMapSize
		offset = (offset + i) % hashMapSize
		return index
	}
}

func getHashValue(data []byte, hashFunc hash.Hash) uint {
	hash := hashFunc.Sum(data)
	return uint(binary.LittleEndian.Uint32(hash[:HASH_CLAMP_VALUE]))
}

func calculateFilterSize(filterSize uint, fpRate float64) uint {
	return uint((math.Ceil(-1 * float64(filterSize) * math.Log(fpRate) / math.Pow(math.Log(2), 2))))
}

func calculateNumberOfHashes(filterSize uint, numberOfElements uint) uint {
	return uint((math.Ceil(float64(filterSize) / float64(numberOfElements) * math.Log(2))))
}
