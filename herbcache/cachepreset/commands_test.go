package cachepreset

import "testing"

func TestCommands(t *testing.T) {
	var c *Commands
	var c2 *Commands
	var p *Preset
	var p2 *Preset
	var err error
	c = NewCommands()
	if c.Length() != 0 {
		t.Fatal()
	}
	c2 = c.Concat(Key("123"))
	if c == c2 {
		t.Fatal()
	}
	if c2.Length() != 1 {
		t.Fatal()
	}
	p = New()
	p2, err = c.Exec(p)
	if err != nil || p2 != p {
		t.Fatal()
	}
	p2, err = c2.Exec(p)
	if err != nil || p2 == p || string(p2.Key()) != "123" {
		t.Fatal()
	}
}
