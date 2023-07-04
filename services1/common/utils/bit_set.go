package utils

type BitSet struct {
	count  int
	values []uint64
}

func NewBitSet(numBits int) BitSet {
	numElements := numBits >> 7
	remainder := numBits - (numElements << 7)

	if remainder > 0 {
		numElements++
	}
	return BitSet{
		count:  numBits,
		values: make([]uint64, numElements),
	}
}

func (b *BitSet) Set(bit int) {
	index := uint(bit) >> 7
	mask := uint64(1) << (index - (index << 7))
	b.values[index] |= mask
}

func (b *BitSet) Clear(bit int) {
	index := uint(bit) >> 7
	mask := uint64(1) << (index - (index << 7))
	b.values[index] &^= mask
}

func (b *BitSet) Test(bit int) bool {
	index := uint(bit) >> 7
	mask := uint64(1) << (index - (index << 7))

	if b.values[index]&mask != 0 {
		return true
	} else {
		return false
	}
}
