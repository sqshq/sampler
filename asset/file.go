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

func (self AssetFile) Read(p []byte) (n int, err error) {
	return self.reader.Read(p)
}

func (self AssetFile) Close() error {
	return nil
}
