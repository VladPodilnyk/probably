// Package filters implements probabilistic data structures
// such as Bloom filter.
package filters

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"hash"
	"math"

	"github.com/vladpodilnyk/probably/internal/bitarray"
)

const hashClampValue = 4

// A BloomFilter holds a bit array that contains information
// about the presence of an element in the set.
// Plus it holds some additional metadata about the filter.
type BloomFilter struct {
	bitarray       *bitarray.BitArray // underlying bit array
	numberOfBits   uint               // number of bits that is needed to keep the given number of elements given the false positive rate (fpRate)
	numberOfHashes uint               // number of hash functions that is used by the filter
	filterSize     uint               // number of elements that the filter can hold (provided by the user)
	fpRate         float64            // false positive rate (provided by the user)
}

// NewBloomFilter creates a new BloomFilter with the given size (size of the set) and false positive rate (fpRate).
func NewBloomFilter(size uint, fpRate float64) BloomFilter {
	array_bits := calculateFilterSize(size, fpRate)
	array := bitarray.New(array_bits)
	hashes := calculateNumberOfHashes(array_bits, size)
	return BloomFilter{
		bitarray:       array,
		numberOfBits:   array_bits,
		numberOfHashes: hashes,
		filterSize:     size,
		fpRate:         fpRate,
	}
}

// Add adds the given data to the filter.
// Data should be presented as a byte array.
func (filter BloomFilter) Add(data []byte) {
	calculateHash := hashGenerator(data, filter.numberOfBits)
	for i := uint(0); i < filter.numberOfHashes; i++ {
		index := calculateHash(i)
		filter.bitarray.Set(index)
	}
}

// Contains checks if the given data is present in the filter.
// In the same fashion as Add, data should be presented as a byte array.
// Keep in mind that Bloom filter can definitely tell if the element is not present in the set,
// but it can only tell that the element is probably present in the set, meaning you can get false positives.
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

// Merge joins two Bloom filters together and stores the result in the first filter.
// To be able to merge two filters they must have the same configuration (size and fpRate).
func (filter BloomFilter) Merge(value BloomFilter) error {
	// Union and merges of Bloom filters described here
	// https://www.cs.utexas.edu/users/lam/396m/slides/Bloom_filters.pdf
	if filter.filterSize != value.filterSize || filter.fpRate != value.fpRate {
		return errors.New("BloomFilters must have the same configuration (size and fpRate)")
	}
	filter.bitarray.Merge(value.bitarray)
	return nil
}

// Union joins two Bloom filters together and returns the result as a new filter.
// To be able to merge two filters they must have the same configuration (size and fpRate).
// If union operation fails, an error is returned.
func (filter BloomFilter) Union(value BloomFilter) (BloomFilter, error) {
	// Union and merges of Bloom filters described here
	// https://www.cs.utexas.edu/users/lam/396m/slides/Bloom_filters.pdf
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

// Size returns the size of the filter (size of the set).
// This value is provided by the user when creating the filter.
func (filter BloomFilter) Size() uint {
	return filter.filterSize
}

// Clear clears the set by resetting the bit array.
func (filter BloomFilter) Clear() {
	filter.bitarray.Clear()
}

// This generator uses two hash functions (MD5 and SHA1) to generate k hash functions
// hashMapSize is the size of the underlying bit array, since the number of bits
// represents the possible address space.
// The idea is very well described in the following paper
// https://www.eecs.harvard.edu/~michaelm/postscripts/tr-02-05.pdf
func hashGenerator(value []byte, hashMapSize uint) func(uint) uint {
	index := getHashValue(value, md5.New())
	offset := getHashValue(value, sha1.New())

	return func(i uint) uint {
		index = (index + offset) % hashMapSize
		offset = (offset + i) % hashMapSize
		return index
	}
}

// This function uses a hashFunc to generate a hash value for the given data.
// Since the hash values are longer than 64 bits, we clamp the value and take only
// the lower 4 bytes (32 bits).
func getHashValue(data []byte, hashFunc hash.Hash) uint {
	hash := hashFunc.Sum(data)
	return uint(binary.LittleEndian.Uint32(hash[:hashClampValue]))
}

// This function calculates the size of the underlying bit array.
// A good place to learn more about the math behidn this is
// https://en.wikipedia.org/wiki/Bloom_filter#Optimal_number_of_hash_functions
func calculateFilterSize(filterSize uint, fpRate float64) uint {
	return uint((math.Ceil(-1 * float64(filterSize) * math.Log(fpRate) / math.Pow(math.Log(2), 2))))
}

// This function calculates the number of hash functions that is needed to generate hash values.
// A good place to learn more about the math behidn this is
// https://en.wikipedia.org/wiki/Bloom_filter#Optimal_number_of_hash_functions
func calculateNumberOfHashes(numberOfElements uint, filterSize uint) uint {
	return uint((math.Ceil(float64(numberOfElements) / float64(filterSize) * math.Log(2))))
}
