package herbcache

import (
	"bytes"
	"testing"
)

func TestToken(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	_, err := WriteNamespace(buf, []byte("test"))
	if err != nil {
		panic(err)
	}
	ns := buf.Bytes()
	if len(ns) == 4 {
		t.Fatal()
	}
	if !bytes.Equal(ns[0:1], TokenBeforeNamespace) {
		t.Fatal()
	}
	var p *Position
	d := p.Append([]byte("g"), []byte("n")).Append([]byte("g2"), []byte("n2")).RootDirectory()
	buf = bytes.NewBuffer(nil)
	_, err = WriteDirectory(buf, d)
	if err != nil {
		panic(err)
	}
	ds := buf.Bytes()
	if len(ds) != 0 {
		t.Fatal()
	}
	buf = bytes.NewBuffer(ds)
	_, err = WriteDirectory(buf, d.Next)
	if err != nil {
		panic(err)
	}
	ds = buf.Bytes()
	if len(ds) == 0 {
		t.Fatal()
	}
	if !bytes.Equal(ds[0:1], TokenBeforeGroup) || bytes.Equal(ds[0:1], ns[0:1]) {
		t.Fatal()
	}

	buf = bytes.NewBuffer(nil)
	_, err = WriteDirectory(buf, d.Next.Next)
	if err != nil {
		panic(err)
	}
	ds2 := buf.Bytes()
	if len(ds) == 0 {
		t.Fatal()
	}
	if !bytes.Equal(ds2[0:1], TokenBeforeGroup) || bytes.Equal(ds2[0:1], ns[0:1]) {
		t.Fatal()
	}

	buf = bytes.NewBuffer(nil)
	_, err = WriteGroupedKey(buf, []byte("g"), []byte("n"))
	if err != nil {
		panic(err)
	}
	ks := buf.Bytes()
	if len(ks) == 0 {
		t.Fatal()
	}
	if !bytes.Equal(ks[0:1], TokenBeforeGroup) {
		t.Fatal()
	}
	if bytes.Equal(ks, ds) {
		t.Fatal()
	}
	buf = bytes.NewBuffer(nil)
	_, err = WriteDirectories(buf, d)
	if err != nil {
		panic(err)
	}
	da := buf.Bytes()
	if !bytes.Equal(da, append(append([]byte{}, ds...), ds2...)) {
		t.Fatal()
	}
}
