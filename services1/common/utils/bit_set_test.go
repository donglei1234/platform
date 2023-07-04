package utils

import (
	"testing"
)

func TestBitSet(t *testing.T) {
	// TODO: Jack: Task #192 to follow-up on possible BitSet issue
	bitSet := NewBitSet(4)
	bitSet.Clear(0)
	bitSet.Set(0)
	bitSet.Test(0)
	bitSet.Clear(0)
	bitSet.Test(0)
}
