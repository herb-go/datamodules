package ncache

import (
	"io"

	"github.com/herb-go/herbdata/datautil"
)

var tokenBeforeNamespace = []byte{':'}
var tokenBeforeValue = []byte("#")
var tokenBeforePath = []byte{0}
var tokenBeforePrefix = []byte{'/'}

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
