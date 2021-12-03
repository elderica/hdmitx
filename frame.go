package main

const (
	MaxFileSize            = 0xff_ff_ff // fly_by_wireプロトコルで送信可能な最大サイズ
	MaxRawNumbering        = 0xf_ff     // numberingのdatumペイロード部の最大値
	MaxPayloadSizePerFrame = 63         // 画像単位ひとつで運べるペイロードサイズ
)

// Frame は1080ビット(135バイト)から成る画像単位である。
// 一つの画像単位で3バイト*42個/2=63バイトのペイロードを運べる。
type Frame [135]byte

// MakeFrames はバイナリデータから画像単位を作る関数である。
func MakeFrames(bindata []byte) []Frame {
	// TODO: 2つ以上の画像単位に対応する。
	frames := make([]Frame, 1)
	frame := frames[0][:]

	// プリアンブルと予約領域を書きこむ。
	frame[0] = 0b1010_1010
	frame[1] = 0b1000_0000
	frame[2] = 0b0000_0000

	return frames
}

func EncodeNumbering(rawNumbering uint16) []byte {
	rawNumbering %= MaxRawNumbering
	datum := DatumEncodeWord(rawNumbering)
	numbering := make([]byte, 3)
	numbering[0] = byte((datum & 0x00ff0000) >> 16)
	numbering[1] = byte((datum & 0x0000ff00) >> 16)
	numbering[2] = byte((datum & 0x000000ff) >> 16)
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
