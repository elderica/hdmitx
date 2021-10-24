package main

func golayEncode(in []byte, out []byte) {
	n := len(in)
	pin := 0
	pout := 0
	g3 := make([]byte, 3)
	g6 := make([]byte, 6)

	for n >= 3 {
		copy(g3, in[pin:pin+3])
		doEncode(g3, g6)
		copy(out[pout:pout+6], g6)
		pin += 3
		pout += 6
		n -= 3
	}
}

func doEncode(g3 []byte, g6 []byte) {
	var (
		parity     uint16
		frontInfo  uint16
		codedFront uint32
		backInfo   uint16
		codedBack  uint32
	)
	frontInfo = uint16(g3[0])<<4 | uint16(g3[1])>>4
	parity = golayParityTable[frontInfo]
	codedFront = uint32(frontInfo)<<12 | uint32(parity)<<1
	if needBalance(codedFront) {
		codedFront |= 1
	}

	backInfo = uint16(g3[1]&0x0f)<<4 | uint16(g3[2])
	parity = golayParityTable[backInfo]
	codedBack = uint32(backInfo)<<12 | uint32(parity)<<1
	if needBalance(codedBack) {
		codedBack |= 1
	}

	g6[0] = byte((codedFront & 0x00ff0000) >> 16)
	g6[1] = byte((codedFront & 0x0000ff00) >> 8)
	g6[2] = byte(codedFront & 0x000000ff)
	g6[0] = byte((codedBack & 0x00ff0000) >> 16)
	g6[1] = byte((codedBack & 0x0000ff00) >> 8)
	g6[2] = byte(codedBack & 0x000000ff)
}

// needBalance は(codeword & 0x00fffffe)の1の個数を数え、
// 0の個数より少なかったらtrueを返します。
func needBalance(codeword uint32) bool {
	count := 0
	for i := 1; i < 24; i++ {
		if codeword&(1<<i) == 1 {
			count++
		}
	}
	return count <= 12
}
