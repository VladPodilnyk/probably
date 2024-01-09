package bitarray

import (
	"testing"
)

func TestAllocateAppropriateAmoutOfMemory(t *testing.T) {
	// allocate bit arrays
	array1 := New(21)
	array2 := New(32)
	array3 := New(4)

	verifyBitArrayLength(t, array1, 3)
	verifyBitArrayLength(t, array2, 4)
	verifyBitArrayLength(t, array3, 1)
}

func TestSetBitInBitArray(t *testing.T) {
	array1 := New(4)
	array2 := New(10)

	array1.Set(2)
	array2.Set(9)

	if !array1.IsSet(2) || !array2.IsSet(9) {
		t.Errorf("Bit set operation failed")
	}
}

func TestMergeTwoBitArrays(t *testing.T) {
	array1 := New(4)
	array2 := New(4)

	array1.Set(2)
	array2.Set(3)

	array1.Merge(array2)

	if !array1.IsSet(2) || !array1.IsSet(3) {
		t.Errorf("Bit merge operation failed")
	}
}

func TestUnionOfTwoBitArrays(t *testing.T) {
	array1 := New(4)
	array2 := New(4)

	array1.Set(2)
	array2.Set(3)

	result := array1.Union(array2)

	if !result.IsSet(2) || !result.IsSet(3) {
		t.Errorf("Bit union operation failed")
	}
}

func TestClearBitArrayState(t *testing.T) {
	array1 := New(4)
	array1.Set(2)
	array1.Set(3)

	array1.Clear()

	if array1.IsSet(2) {
		t.Errorf("Bit clear operation failed")
	}
}

func verifyBitArrayLength(t *testing.T, array *BitArray, expectedLenth int) {
	if len(array.data) != expectedLenth {
		t.Errorf("Expected array size to be %d, got %d", expectedLenth, len(array.data))
	}
}
