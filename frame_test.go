package main

import (
	"testing"
)

func TestSplitByteUpto(t *testing.T) {
	tests := []struct {
		input         []byte
		n             uint
		expectedFront []byte
		expectedBack  []byte
	}{
		{
			[]byte{},
			0,
			[]byte{},
			[]byte{},
		},
		{
			[]byte{1, 2, 3},
			0,
			[]byte{},
			[]byte{1, 2, 3},
		},
		{
			[]byte{1, 2, 3},
			4,
			[]byte{1, 2, 3},
			[]byte{},
		},
		{
			[]byte{1, 2, 3},
			3,
			[]byte{1, 2, 3},
			[]byte{},
		},
		{
			[]byte{1, 2, 3},
			2,
			[]byte{1, 2},
			[]byte{3},
		},
	}

	for _, test := range tests {
		gotFront, gotBack := SplitByteUpto(test.input, test.n)
		if !compareBytes(gotFront, test.expectedFront) {
			t.Errorf("TakeByteUpto(%v, %d) = (%v ,%v), want %v",
				test.input, test.n, gotFront, gotBack, test.expectedFront)
		}
	}
}
