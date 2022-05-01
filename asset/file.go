package asset

import (
	"bytes"
	"io"
)

type File struct {
	reader io.Reader
}

func NewAssetFile(data []byte) File {
	return File{bytes.NewReader(data)}
}

func (a File) Read(p []byte) (n int, err error) {
	return a.reader.Read(p)
}

func (a File) Close() error {
	return nil
}
