package herbcache

import (
	"io"

	"github.com/herb-go/herbdata/datautil"
)

var TokenBeforeNamespace = []byte{':'}
var TokenBeforeKey = []byte("#")
var TokenBeforeDirectory = []byte{0}
var TokenBeforeGroup = []byte{'/'}

func writeTokenAndData(w io.Writer, token []byte, data []byte) (int, error) {
	var written int
	var l int
	var err error
	l, err = w.Write(token)
	if err != nil {
		return 0, err
	}
	written += l
	l, err = datautil.WriteLengthBytes(w, len(data))
	if err != nil {
		return 0, err
	}
	written += l
	l, err = w.Write(data)
	if err != nil {
		return 0, err
	}
	written += l
	return written, nil
}

func WriteNamespace(w io.Writer, data []byte) (int, error) {
	return writeTokenAndData(w, TokenBeforeNamespace, data)
}

func WriteDirectory(w io.Writer, d *Directory) (int, error) {
	pl, err := writeTokenAndData(w, TokenBeforeGroup, d.Group)
	if err != nil {
		return 0, nil
	}
	dl, err := writeTokenAndData(w, TokenBeforeDirectory, d.Name)
	if err != nil {
		return 0, nil
	}
	return pl + dl, nil
}
func WriteGroupedKey(w io.Writer, group []byte, key []byte) (int, error) {
	pl, err := writeTokenAndData(w, TokenBeforeGroup, group)
	if err != nil {
		return 0, nil
	}
	kl, err := writeTokenAndData(w, TokenBeforeDirectory, key)
	if err != nil {
		return 0, nil
	}
	return pl + kl, nil
}
