package main

import (
	"testing"
)

func TestEnsureBalance(t *testing.T) {
	tests := []struct {
		input    uint32
		expected uint32
	}{
		{0b0000_0000_0000_0000_0000_0000, 0b0000_0000_0000_0000_0000_0001},
		{0b1111_1111_1111_1111_1111_1110, 0b1111_1111_1111_1111_1111_1110},
		{0b1010_1010_1010_1010_1010_1010, 0b1010_1010_1010_1010_1010_1010},
		{0b1010_1010_1010_1010_1010_1011, 0b1010_1010_1010_1010_1010_1011},
		{0b1010_1010_1010_1010_1010_1000, 0b1010_1010_1010_1010_1010_1001},
		{0b0010_0010_0010_0010_0010_0010, 0b0010_0010_0010_0010_0010_0011},
	}

	for _, test := range tests {
		got := EnsureBalance(test.input)
		if got != test.expected {
			t.Errorf("EnsureBalance(%08x) = %08x, want %08x", test.input, got, test.expected)
		}
	}
}

func TestDatumEncode(t *testing.T) {
	tests := []struct {
		input    []byte
		expected []byte
	}{
		{
			[]byte{},
			[]byte{},
		},
		{
			[]byte{0, 0, 0},
			[]byte{0, 0, 1, 0, 0, 1},
		},
		{
			[]byte{0, 0},
			[]byte{0, 0, 1, 0, 0, 1},
		},
		{
			[]byte{0},
			[]byte{0, 0, 1, 0, 0, 1},
		},
		{
			[]byte{0, 0, 0, 0, 0},
			[]byte{
				0, 0, 1, 0, 0, 1,
				0, 0, 1, 0, 0, 1,
			},
		},
		{
			[]byte{0, 0, 0, 0},
			[]byte{
				0, 0, 1, 0, 0, 1,
				0, 0, 1, 0, 0, 1,
			},
		},
		{
			[]byte{1, 0, 2},
			[]byte{
				0x75, 0x0c, 0x01,
				0x9f, 0x14, 0x01,
			},
		},
	}

	for _, test := range tests {
		got := DatumEncode(test.input)
		if !compareBytes(got, test.expected) {
			t.Errorf("EnsureBalance(%v) = %v, want %v",
				test.input, got, test.expected)
		}
	}
}

// compareBytes は a と b が同じ内容であるか確認する。
func compareBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
