package main

import (
	"log"
	"os"
)

const (
	MaxSize = 0x00ffffff // fly_by_wireプロトコルで送信可能な最大サイズ
)

type Frame []byte

type TooLargeSizeError struct{}

func (err *TooLargeSizeError) Error() string {
	return "16777215(0x00ffffff)バイトより大きなサイズのデータです"
}

// constructFrame は画像単位を生成する関数です。
func constructFrame(size []byte, numbering []byte, payload []byte) (frame Frame, nconsumed uint) {
	frame = make(Frame, 1920)
	frame[0] = 0b10101010
	frame[1] = 0b10000111
	frame[2] = 0b10000101
	copy(frame[3:6], size)
	copy(frame[6:9], numbering)

	// 一つの画像単位にはdatumは42個まで入れらえる。
	nDatumToMake := len(payload) / 3
	if nDatumToMake > 42 {
		nDatumToMake = 42
	}
	for n := 0; n < nDatumToMake; n++ {
		copy(frame[9+3*n:9+3*n+3], payload[3*n:3*n+3])
	}
}

func loadFile(filepath string) ([]byte, error) {
	var (
		fileinfo os.FileInfo
		err      error
	)
	fileinfo, err = os.Stat(filepath)
	if err != nil {
		log.Printf("cannot stat file: %s\n", err)
		return nil, err
	}
	if fileinfo.Size() > MaxSize {
		return nil, &TooLargeSizeError{}
	}

	contents, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

// constructFramesFromFile は filepath のファイルを読み取って画像単位群を生成する。
func constructFramesFromFile(filepath string) ([]Frame, error) {
	contents, err := loadFile(filepath)
	if err != nil {
		return nil, err
	}

	// golayEncodeに渡すスライスは3の倍数の長さでなくてはならない。
	padsize := len(contents) % 3
	if padsize == 1 {
		contents = append(contents, 0)
	} else if padsize == 2 {
		contents = append(contents, 0, 0)
	}

	size := len(contents) / 3
	sizebytes := []byte{
		byte((size & 0x00ff0000) >> 16),
		byte((size & 0x0000ff00) >> 8),
		byte((size & 0x000000ff)),
	}

	payload := make([]byte, len(contents)*2)
	golayEncode(contents, payload)

	frames := make([]Frame, 0)
	numbering := uint32(0)
	begin := uint(0)
	for begin < uint(len(contents)) {
		numberingbytes := []byte{
			byte((numbering & 0x00ff0000) >> 16),
			byte((numbering & 0x0000ff00) >> 8),
			byte((numbering & 0x000000ff)),
		}
		frame, nconsumed := constructFrame(sizebytes, numberingbytes, payload[begin:])
		begin += nconsumed
		frames = append(frames, frame)
		numbering++
	}

	return frames, nil
}
