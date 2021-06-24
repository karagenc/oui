package ouidata

import (
	"bytes"
	"compress/gzip"
	_ "embed"

	"github.com/tomruk/oui"
)

//go:embed oui.txt.gz
var data []byte

func NewDB() (*oui.DB, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return oui.NewDBFromReader(reader)
}
