package cachepreset

import (
	"bytes"
	"testing"

	"github.com/herb-go/datamodules/herbcache"
	"github.com/herb-go/herbdata/dataencoding"
)

func TestOverride(t *testing.T) {
	var p *Preset
	var p2 *Preset
	var err error
	p = New()
	p2 = p.OverrideLockers(NewLockers())
	if p == p2 || p2.Lockers() == nil {
		t.Fatal()
	}
	p2 = p.OverrideKey([]byte("key"))
	if p == p2 || string(p2.Key()) != "key" {
		t.Fatal()
	}
	p2 = p.OverrideLoader(Loader(func([]byte) ([]byte, error) {
		return nil, nil
	}))
	if p == p2 || p2.Loader() == nil {
		t.Fatal()
	}
	p2 = p.OverrideData([]byte("data"))
	if p == p2 || string(p2.Data()) != "data" {
		t.Fatal()
	}
	p2 = p.OverrideEncoding(dataencoding.NopEncoding)
	if p == p2 || p2.Encoding() != dataencoding.NopEncoding {
		t.Fatal()
	}
	p2 = p.OverrideTTL(15)
	if p == p2 || p2.TTL() != 15 {
		t.Fatal()
	}
	cache := herbcache.New()
	p, err = Cache(cache).Exec(p)
	if err != nil {
		t.Fatal()
	}
	p2, err = p.Flushable(true).Apply()
	if p == p2 || err != nil || p2.Cache().Flushable() == false {
		t.Fatal()
	}
	p2, err = p.Allocate("namespace").Apply()
	if p == p2 || err != nil || string(p2.Cache().Namespace()) != "namespace" {
		t.Fatal()
	}
	p2, err = p.ChildCache("child").Apply()
	if p == p2 || err != nil || string(p2.Cache().Position().Name) != "child" {
		t.Fatal()
	}
	p2, err = p.PrefixCache("prefix").Apply()
	if p == p2 || err != nil || string(p2.Cache().Group()) != "prefix" {
		t.Fatal()
	}
}

func TestClone(t *testing.T) {
	p := New()
	p2 := p.Clone()
	if p == p2 {
		t.Fatal()
	}
	p.data = []byte("12345")
	if bytes.Compare(p.data, p2.data) == 0 {
		t.Fatal()
	}
}
