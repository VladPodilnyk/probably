package filters

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/binary"
	"math"

	"github.com/probably/internal/bitarray"
)

const HASH_CLAMP_VALUE = 4

type BloomFilter struct {
	bitarray       *bitarray.BitArray
	numberOfBits   uint
	numberOfHashes uint
}

func NewBloomFilter(size uint, fpRate float64) BloomFilter {
	array_bits := calculateFilterSize(size, fpRate)
	array := bitarray.New(array_bits)
	hashes := calculateNumberOfHashes(size, array_bits)
	return BloomFilter{bitarray: array, numberOfBits: array_bits, numberOfHashes: hashes}
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

func (filter BloomFilter) Merge(value BloomFilter) {}

func (filter BloomFilter) Union(value BloomFilter) {}

func (filter BloomFilter) Clear() {
	filter.bitarray.Clear()
}

func hashGenerator(value []byte, hashMapSize uint) func(uint) uint {
	index := getMd5Hash(value)
	offset := getSha1Hash(value)

	return func(i uint) uint {
		index = (index + offset) % hashMapSize
		offset = (offset + i) % hashMapSize
		return index
	}
}

// TODO: abstract this into a common functino that does clamping
func getMd5Hash(value []byte) uint {
	fullHash := md5.Sum(value)
	return uint(binary.LittleEndian.Uint32(fullHash[:HASH_CLAMP_VALUE]))
}

func getSha1Hash(value []byte) uint {
	fullHash := sha1.Sum(value)
	return uint(binary.LittleEndian.Uint32(fullHash[:HASH_CLAMP_VALUE]))
}

func calculateFilterSize(filterSize uint, fpRate float64) uint {
	return uint((math.Ceil(-1 * float64(filterSize) * math.Log(fpRate) / math.Pow(math.Log(2), 2))))
}

func calculateNumberOfHashes(filterSize uint, numberOfElements uint) uint {
	return uint((math.Ceil(float64(filterSize) / float64(numberOfElements) * math.Log(2))))
}
