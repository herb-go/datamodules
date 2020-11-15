package rcache

import (
	"bytes"
	"io"

	"github.com/herb-go/herbdata/datautil"
)

type enity struct {
	typecode byte
	version  []byte
	data     []byte
}

func (e *enity) SaveTo(w io.Writer) error {
	var err error
	_, err = w.Write([]byte{e.typecode})
	if err != nil {
		return err
	}
	if e.typecode == enityTypecodeRevocable {
		err = datautil.PackTo(w, nil, e.version)
		if err != nil {
			return err
		}
	}
	_, err = w.Write(e.data)
	return err
}

const enityTypecodeRevocable = byte(0)
const enityTypecodeIrrevocable = byte(1)

func newEnity() *enity {
	return &enity{}
}
func createEnity(revocable bool, version []byte, data []byte) *enity {
	e := newEnity()
	if revocable {
		e.typecode = enityTypecodeRevocable
		e.version = version
	} else {
		e.typecode = enityTypecodeIrrevocable
	}
	e.data = data
	return e
}

func loadEnity(data []byte, revocable bool, version []byte) (*enity, error) {
	datalength := len(data)
	if datalength == 0 {
		return nil, ErrUnresolvableCacheEnity
	}
	switch data[0] {
	case enityTypecodeIrrevocable:
		if revocable {
			return nil, ErrEnityTypecodeNotMatch
		}
		return createEnity(false, nil, data[1:]), nil
	case enityTypecodeRevocable:

		if !revocable {
			return nil, ErrEnityTypecodeNotMatch
		}
		buf := bytes.NewBuffer(data[1:])
		versiondata, err := datautil.UnpackFrom(buf, nil)
		if err != nil {
			return nil, err
		}
		if bytes.Compare(version, versiondata) != 0 {
			return nil, ErrEnityVersionNotMatch
		}
		return createEnity(false, versiondata, buf.Bytes()), nil
	}
	return nil, ErrUnresolvableCacheEnity
}
