package main

const (
	MaxPayloadSize  = 0xff_ff_ff // fly_by_wireプロトコルで送信可能な最大サイズ
	MaxRawNumbering = 0xf_ff     // numberingのdatumペイロード部の最大値

	MaxNDatum              = 42              // 画像単位ひとつに載せられるdatumの個数
	MaxDataSize            = 3 * MaxNDatum   // 画像単位のdataの長さ
	MaxPayloadSizePerFrame = MaxDataSize / 2 // 画像単位ひとつで運べるペイロードサイズ
)

// Frame は1080ビット(135バイト)から成る画像単位である。
// 一つの画像単位で3バイト*42個/2=63バイトのペイロードを運べる。
type Frame [135]byte

// MakeFrames はバイナリデータから画像単位を作る関数である。
func MakeFrames(bindata []byte) []Frame {
	data := DatumEncode(bindata)

	frames := make([]Frame, 0, 10)
	for fi := 0; len(data) > 0; fi++ {
		frames = append(frames, Frame{})
		frame := frames[fi][:]

		// プリアンブルと予約領域を書きこむ。
		frame[0] = 0b1010_1010
		frame[1] = 0b1000_0000
		frame[2] = 0b0000_0000

		// size(24ビット幅)を書きこむ。
		size := uint32(len(bindata))
		frame[3] = byte((size & 0x00ff0000) >> 16)
		frame[4] = byte((size & 0x0000ff00) >> 8)
		frame[5] = byte(size & 0x000000ff)

		// numbering(24ビット幅)を書きこむ。
		rawNumbering := uint16(0)
		numbering := EncodeNumbering(rawNumbering)
		copy(frame[6:9], numbering)

		dataFront, dataBack := SplitByteUpto(data, MaxDataSize)
		copy(frame[9:], dataFront)

		data = dataBack
	}

	return frames
}

func EncodeNumbering(rawNumbering uint16) []byte {
	rawNumbering %= MaxRawNumbering + 1
	datum := DatumEncodeWord(rawNumbering)
	numbering := make([]byte, 3)
	numbering[0] = byte((datum & 0x00ff0000) >> 16)
	numbering[1] = byte((datum & 0x0000ff00) >> 8)
	numbering[2] = byte(datum & 0x000000ff)
	return numbering
}

// SplitByteUpto は、sliceの先頭からn個を境界として2つに分割する。
// もしlen(slice) < nならば、sliceと空スライスを返す。
func SplitByteUpto(slice []byte, n uint) ([]byte, []byte) {
	if len(slice) < int(n) {
		return slice, []byte{}
	}

	return slice[0:n], slice[n:]
}
