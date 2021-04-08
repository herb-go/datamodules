package herbcache

import (
	"bytes"
	"testing"
)

func newTestStorage() *Storage {
	s := NewStorage()
	return s
}
func TestCacheOperations(t *testing.T) {
	var cc *Cache
	s := newTestStorage()
	c := New().OverrideStorage(s).OverrideFlushable(true)
	if c.Storage() != s {
		t.Fatal(c)
	}
	if c.Equal(nil) {
		t.Fatal(c)
	}
	if !cc.Equal(nil) {
		t.Fatal(cc)
	}
	cc = c.Clone()
	if !cc.Equal(cc) {
		t.Fatal(cc)
	}

	cc = c.OverrideFlushable(false)
	if cc.Flushable() != false || c.Flushable() != true {
		t.Fatal(cc, c)
	}
	cc = c.OverrideStorage(nil)
	if cc.Storage() != nil || c.Storage() == nil {
		t.Fatal(cc, c)
	}
	cc = c.OverrideGroup([]byte("group"))
	if bytes.Equal(c.Group(), cc.Group()) {
		t.Fatal(cc, c)
	}
	if cc.Equal(c) {
		t.Fatal(cc, c)
	}
	cc = c.Migrate([]byte("newns"))
	if bytes.Equal(cc.Namespace(), c.Namespace()) {
		t.Fatal(cc, c)
	}
	if cc.Equal(c) {
		t.Fatal(cc, c)
	}
	if c.Position() != nil {
		t.Fatal()
	}
	cc = c.Migrate([]byte("newns")).SubCache([]byte("child")).OverrideGroup([]byte("group")).Migrate(nil)
	if c.Position() != nil || c.Namespace() != nil || c.Group() != nil {
		t.Fatal()
	}
	cc = c.SubCache([]byte("child")).OverrideGroup([]byte("group"))
	if cc.Equal(c) {
		t.Fatal(cc, c)
	}
	if cc.Position() == nil {
		t.Fatal()
	}
	cc = c.OverrideStorage(nil).OverrideFlushable(false)
	c.CopyFrom(cc)
	if c.flushable != false ||
		c.storage != nil {
		t.Fatal(c)
	}
}

func TestStringCache(t *testing.T) {
	c := New()
	if !c.OverrideGroup([]byte("g")).Equal(c.PrefixCache("g")) {
		t.Fatal()
	}
	if !c.SubCache([]byte("sub")).Equal(c.ChildCache("sub")) {
		t.Fatal()
	}
	if !c.Allocate("ns").Equal(c.Migrate([]byte("ns"))) {
		t.Fatal()
	}
}

func TestSetCache(t *testing.T) {
	c := New()
	dst := New()
	if c == dst {
		t.Fatal()
	}
	dst = dst.OverrideFlushable(true)
	SetCache(c, dst)
	if c == dst {
		t.Fatal()
	}
	if !c.Equal(dst) {
		t.Fatal()
	}
}
