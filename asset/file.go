package asset

import (
	"bytes"
	"io"
)

type AssetFile struct {
	reader io.Reader
}

func NewAssetFile(data []byte) AssetFile {
	return AssetFile{bytes.NewReader(data)}
}

func (a AssetFile) Read(p []byte) (n int, err error) {
	return a.reader.Read(p)
}

func (a AssetFile) Close() error {
	return nil
}
