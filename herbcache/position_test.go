package herbcache

import (
	"bytes"
	"testing"
)

func TestPosition(t *testing.T) {
	var p *Position
	if p.Equal(new(Position)) {
		t.Fatal()
	}
	p2 := p.Append([]byte("g2"), []byte("n2"))
	if p2.Equal(p) {
		t.Fatal(p, p2)
	}
	if !p2.Equal(p.Append([]byte("g2"), []byte("n2"))) {
		t.Fatal()
	}
	if p2.Equal(p.Append([]byte("g2"), []byte("n3"))) {
		t.Fatal()
	}
	if p.Append(nil, []byte("n1")).Append([]byte("g2"), []byte("n2")).
		Equal(
			p.Append([]byte("g1"), []byte("n1")).Append([]byte("g2"), []byte("n2")),
		) {
		t.Fatal()
	}
}

func TestDirectory(t *testing.T) {
	var p *Position
	p = p.Append([]byte("g1"), []byte("n1")).Append([]byte("g2"), []byte("n2"))
	d := p.RootDirectory()
	if d == nil {
		t.Fatal()
	}
	d2 := d.Next
	if d2 == nil || !bytes.Equal(d2.Group, []byte("g1")) || !bytes.Equal(d2.Name, []byte("n1")) {
		t.Fatal()
	}
	d3 := d2.Next
	if d3 == nil || !bytes.Equal(d3.Group, []byte("g2")) || !bytes.Equal(d3.Name, []byte("n2")) || d3.Next != nil {
		t.Fatal()
	}
}
