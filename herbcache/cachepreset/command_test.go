package cachepreset

import (
	"testing"

	"github.com/herb-go/datamodules/herbcache"
	"github.com/herb-go/herbdata/dataencoding/msgpackencoding"
)

func TestCommand(t *testing.T) {
	var np *Preset
	var cp *Preset
	var err error
	p := New()
	np, err = Key([]byte("12345")).Exec(p)
	if np == p || err != nil || string(np.key) != "12345" {
		t.Fatal()
	}
	np, err = Data([]byte("12345")).Exec(p)
	if np == p || err != nil || string(np.data) != "12345" {
		t.Fatal()
	}
	np, err = TTL(15).Exec(p)
	if np == p || err != nil || np.ttl != 15 {
		t.Fatal()
	}
	loader := func([]byte) ([]byte, error) { return nil, nil }
	np, err = Loader(loader).Exec(p)
	if np == p || err != nil || np.loader == nil {
		t.Fatal()
	}
	lockers := NewLockers()
	np, err = lockers.Exec(p)
	if np == p || err != nil || np.lockers != lockers {
		t.Fatal()
	}

	encoding := Encoding(msgpackencoding.Encoding)
	np, err = encoding.Exec(p)
	if np == p || err != nil || np.encoding != msgpackencoding.Encoding {
		t.Fatal()
	}
	cache := herbcache.New()
	cp, err = Cache(cache).Exec(p)
	if cp == p || err != nil || cp.cache != cache {
		t.Fatal()
	}
	np, err = Flushable(true).Exec(cp)
	if cp == np || err != nil || np.cache.Flushable() != true {
		t.Fatal()
	}
	np, err = Allocate("namespace").Exec(cp)
	if cp == np || err != nil || string(np.cache.Namespace()) != "namespace" {
		t.Fatal()
	}
	np, err = ChildCache("child").Exec(cp)
	if cp == np || err != nil || string(np.cache.Position().Name) != "child" {
		t.Fatal()
	}
	np, err = PrefixCache("prefix").Exec(cp)
	if cp == np || err != nil || string(np.cache.Group()) != "prefix" {
		t.Fatal()
	}
}
