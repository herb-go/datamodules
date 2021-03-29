package herbcache

import "testing"

func TestDirective(t *testing.T) {
	s := newTestStorage()
	c := New()
	cc := New()
	if !c.Equal(cc) {
		t.Fatal()
	}
	c = c.Migrate([]byte("ns")).OverrideFlushable(true).OverrideGroup([]byte("g")).OverrideStorage(s).SubCache([]byte("sub"))
	if c.Equal(cc) {
		t.Fatal()
	}
	if cc.IsPreparing() {
		t.Fatal()
	}
	cc = Prepare(cc, Migrate([]byte("ns")), Flushable(true), Group([]byte("g")), s, SubCache([]byte("sub")))
	if c.Equal(cc) {
		t.Fatal()
	}
	if !cc.IsPreparing() {
		t.Fatal()
	}
	err := cc.Ready()
	if err != nil {
		panic(err)
	}
	if cc.IsPreparing() {
		t.Fatal()
	}
	if !c.Equal(cc) {
		t.Fatal()
	}
}
