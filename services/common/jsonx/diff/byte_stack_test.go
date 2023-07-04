package diff

import (
	"testing"
)

func TestByteStack(t *testing.T) {
	var testCase = []struct {
		arrayOfByteArrays [][]byte
	}{
		// test case 1
		{[][]byte{
			{byte(1), byte(2), byte(3)},
			{byte(4), byte(5), byte(6)},
			{byte(7), byte(8), byte(9)}},
		},

		// test case 2
		{[][]byte{
			{byte('h'), byte('e'), byte('l'), byte('o'), byte('o')},
			{byte('w'), byte('o'), byte('r'), byte('l'), byte('d')}},
		},

		// test case 3
		{[][]byte{}},
	}

	for _, tc := range testCase {
		var capacity = uint(len(tc.arrayOfByteArrays))

		byteStack := NewByteStack(capacity)

		if byteStack.IsEmpty() == false {
			t.Error("Error: Newly instantiated ByteStack did not return expected IsEmpty() value of true.")
		}

		if len(byteStack.Data()) != 0 {
			t.Error("Error: Newly instantiated ByteStack did not return expected Data() len of 0.")
		}

		if byteStack.NumElements() != 0 {
			t.Error("Error: Newly instantiated ByteStack did not return expected NumElements() value of 0.")
		}

		byteStack.Pop()

		if byteStack.IsEmpty() == false {
			t.Error("Error: Newly instantiated ByteStack did not return expected IsEmpty() value of true.")
		}

		if len(byteStack.Data()) != 0 {
			t.Error("Error: Newly instantiated ByteStack did not return expected Data() len of 0.")
		}

		if byteStack.NumElements() != 0 {
			t.Error("Error: Newly instantiated ByteStack did not return expected NumElements() value of 0.")
		}

		byteArraysCount := len(tc.arrayOfByteArrays)

		bytesCount := 0
		for _, arrayOfBytes := range tc.arrayOfByteArrays {
			bytesCount += len(arrayOfBytes)
		}

		for _, arrayOfBytes := range tc.arrayOfByteArrays {
			byteStack.Push(arrayOfBytes)
		}

		if byteStack.NumElements() != byteArraysCount {
			t.Error("Error: ByteStack.NumElements() value did not match expected count. Expected:",
				byteArraysCount, "| Actual:", byteStack.NumElements())
		}

		if bytesCount > 0 && byteStack.IsEmpty() {
			t.Error("Error: ByteStack returned unexpected IsEmpty() value of true after", bytesCount,
				"bytes pushed onto the stack.")
		}

		loopCount := 0
		for byteStack.NumElements() > 0 {
			loopCount++
			byteStack.Pop()

			if loopCount > byteArraysCount {
				t.Error("ByteStack.NumElements() value did not hit 0 despite a Pop() per earlier Push()")
				break
			}
		}

		if byteStack.IsEmpty() == false {
			t.Error("Error: ByteStack did not return expected IsEmpty() value of true after all elements popped")
		}

		if len(byteStack.Data()) != 0 {
			t.Error("Error: ByteStack did not return expected Data() len of 0 after all elements popped")
		}

		if byteStack.NumElements() != 0 {
			t.Error("Error: ByteStack did not return expected NumElements() value of 0 after all elements popped")
		}
	}
}
