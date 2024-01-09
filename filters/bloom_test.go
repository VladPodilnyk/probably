package filters

import (
	"testing"
)

func TestBloomFilterDefinatlyKnowsIfElementIsNotPresent(t *testing.T) {
	filter := NewBloomFilter(100, 0.01)
	if filter.Contains([]byte("hello")) {
		t.Errorf("BloomFilter should not contain the element")
	}
}

// This test might be flaky, but given low false positive rate it should be fine
func TestSucceccfullyAddElementToBloomFilter(t *testing.T) {
	filter := NewBloomFilter(100, 0.001)
	filter.Add([]byte("hello"))
	if !filter.Contains([]byte("hello")) {
		t.Errorf("BloomFilter should contain the element")
	}
}

func TestSuccessfullyMergeTwoBloomFilters(t *testing.T) {
	filter1 := NewBloomFilter(100, 0.001)
	filter2 := NewBloomFilter(100, 0.001)

	filter1.Add([]byte("hello"))
	filter2.Add([]byte("world"))

	filter1.Merge(filter2)

	if !filter1.Contains([]byte("hello")) || !filter1.Contains([]byte("world")) {
		t.Errorf("BloomFilter should contain both elements")
	}
}

func TestSuccessfullUnionOfTwoBloomFilters(t *testing.T) {
	filter1 := NewBloomFilter(100, 0.001)
	filter2 := NewBloomFilter(100, 0.001)

	filter1.Add([]byte("hello"))
	filter2.Add([]byte("world"))

	result, err := filter1.Union(filter2)

	if err != nil {
		t.Errorf("Union operation failed")
	}

	if !result.Contains([]byte("hello")) || !result.Contains([]byte("world")) {
		t.Errorf("BloomFilter should contain both elements")
	}
}

func TestFailedToMergeTwoBloomFiltersDueToDifferentConfiguration(t *testing.T) {
	errors := []error{}
	filter1 := NewBloomFilter(10, 1.00)
	filter2 := NewBloomFilter(100, 0.001)
	filter3 := NewBloomFilter(100, 0.1)
	filter4 := NewBloomFilter(5, 0.001)

	allParamsAreDifferent := filter1.Merge(filter2)
	fpRateIsDifferent := filter2.Merge(filter3)
	sizeIsDifferent := filter2.Merge(filter4)

	errors = append(errors, allParamsAreDifferent, fpRateIsDifferent, sizeIsDifferent)
	if len(errors) != 3 {
		t.Errorf("Some filters with different configurations are merged")
	}
}
