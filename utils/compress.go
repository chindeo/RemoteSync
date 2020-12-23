package utils

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io"
	"os"
)

func Compress(data []byte) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	flateWrite, err := flate.NewWriter(buf, flate.BestCompression)
	if err != nil {
		return nil, err
	}
	defer flateWrite.Close()
	io.Copy(flateWrite, bytes.NewReader(data))

	return buf, nil
}
func UnCompress(data []byte) {
	flateReader := flate.NewReader(bytes.NewReader(data))
	defer flateReader.Close()
	// 输出
	fmt.Println("解压后的内容：")
	io.Copy(os.Stdout, flateReader)
}
