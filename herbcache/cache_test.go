package herbcache

import (
	"bytes"
	"testing"
)

type testStorage struct {
	NopStorage
}

func TestCacheOperations(t *testing.T) {
	var cc *Cache
	s := &testStorage{}
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
	cc = c.SubCache([]byte("child")).OverrideGroup([]byte("group"))
	if cc.Equal(c) {
		t.Fatal(cc, c)
	}
	cc = c.OverrideStorage(nil).OverrideFlushable(false)
	c.CopyFrom(cc)
	if c.flushable != false ||
		c.storage != nil {
		t.Fatal(c)
	}
}
