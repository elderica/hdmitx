package main

/* このプログラムではDatumと呼ばれるフォーマットを採用する。
 * 12ビットのペイロード + 11ビットのパリティ + バランスビット1ビットで構成される。
 * パリティには完全ゴレイ符号を用いる。
 * バランスビットは、ペイロードとパリティに出現したビットのうち、出現数が少ないものを採用する。
 * バランスビットを用いると、受信特性が改善することがある。
 */

// DatumEncode はバイナリデータをDatum列にエンコードする。
// バイナリデータの長さが3の倍数でないとき、パディングを行ってからエンコードする。
func DatumEncode(in []byte) (out []byte) {
	padsize := 3 - len(in)%3
	if padsize == 3 {
		padsize = 0
	}
	n := len(in) + padsize

	outsize := n * 2
	out = make([]byte, outsize)

	paddedIn := make([]byte, n)
	copy(paddedIn, in)

	pin := 0
	pout := 0
	g3 := make([]byte, 3)
	g6 := make([]byte, 6)

	for n >= 3 {
		copy(g3, paddedIn[pin:pin+3])
		doEncode(g3, g6)
		copy(out[pout:pout+6], g6)
		pin += 3
		pout += 6
		n -= 3
	}

	return
}

func doEncode(g3 []byte, g6 []byte) {
	var (
		parity     uint16 // 11ビット
		frontInfo  uint16 // 12ビット
		codedFront uint32 // 24ビット
		backInfo   uint16 // 12ビット
		codedBack  uint32 // 24ビット
	)
	frontInfo = (uint16(g3[0]) << 4) | (uint16(g3[1]) >> 4)
	parity = golayParityTable[frontInfo]
	codedFront = (uint32(frontInfo)<<11 | uint32(parity)) << 1
	codedFront = EnsureBalance(codedFront)

	backInfo = (uint16(g3[1]&0x0f) << 4) | uint16(g3[2])
	parity = golayParityTable[backInfo]
	codedBack = (uint32(backInfo)<<11 | uint32(parity)) << 1
	codedBack = EnsureBalance(codedBack)

	g6[0] = byte((codedFront & 0x00ff0000) >> 16)
	g6[1] = byte((codedFront & 0x0000ff00) >> 8)
	g6[2] = byte(codedFront & 0x000000ff)
	g6[3] = byte((codedBack & 0x00ff0000) >> 16)
	g6[4] = byte((codedBack & 0x0000ff00) >> 8)
	g6[5] = byte(codedBack & 0x000000ff)
}

// DatumEncodeWord は、12ビットのplainwordをDatumに変換する。
func DatumEncodeWord(plainword uint16) (datum uint32) {
	parity := golayParityTable[plainword]
	datum = (uint32(plainword)<<11 | uint32(parity)) << 1
	datum = EnsureBalance(datum)
	return
}

// EnsureBalance は(codeword & 0x00fffffe)の1の個数を数え、
// 0の個数と1の個数が均等に近付くようにバランスビットを付加する。
func EnsureBalance(codeword uint32) uint32 {
	count := 0
	for i := 1; i < 24; i++ {
		if codeword&(1<<i) > 0 {
			count++
		}
	}
	if count < 12 {
		return codeword | 1
	}
	return codeword
}
