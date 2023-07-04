package diff

// ByteStack encapsulates building a concatenation of byte buffers.
type ByteStack struct {
	data    []byte
	lengths []int
}

// Return a new ByteStack with the given initial capacity.
func NewByteStack(capacity uint) *ByteStack {
	return &ByteStack{
		data:    make([]byte, 0, capacity),
		lengths: make([]int, 0, 32),
	}
}

// Push the provided values as a single element on the stack.
func (b *ByteStack) Push(values ...[]byte) {
	length := 0

	for _, value := range values {
		b.data = append(b.data, value...)
		length += len(value)
	}

	b.lengths = append(b.lengths, length)
}

// Pop the last Push operation from the stack.
func (b *ByteStack) Pop() {
	if b.IsEmpty() == false {
		lastLengthIndex := len(b.lengths) - 1
		b.data = b.data[:len(b.data)-b.lengths[lastLengthIndex]]
		b.lengths = b.lengths[:lastLengthIndex]
	}
}

// Get the current stack data.
func (b ByteStack) Data() []byte {
	return b.data
}

// Get the number of elements pushed on the stack.
func (b ByteStack) NumElements() int {
	return len(b.lengths)
}

// Is this ByteStack empty?
func (b ByteStack) IsEmpty() bool {
	return len(b.lengths) == 0
}
