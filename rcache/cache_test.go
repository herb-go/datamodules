package rcache

import (
	"bytes"
	"sync"
	"testing"
	"time"

	"github.com/herb-go/herbdata"

	"github.com/herb-go/herbdata-drivers/kvdb-drivers/freecachedb"
	"github.com/herb-go/herbdata/kvdb/commonkvdb"
)

var _ herbdata.Cache = New()

func newTestCache() *Cache {
	c := New()
	e := NewEngine()
	e.VersionStore = commonkvdb.NewInMemory()
	d, err := (&freecachedb.Config{Size: 50000}).CreateDriver()
	if err != nil {
		panic(err)
	}
	e.Store = d
	return c.WithEngine(e)
}

var TestKey = []byte("testkey")
var TestKey2 = []byte("testkey2")

var TestData = []byte("testdata")

func TestCache(t *testing.T) {
	var err error
	var data []byte
	c := newTestCache()
	if c.Irrevocable() {
		t.Fatal(c.Irrevocable())
	}
	err = c.SetWithTTL(TestKey, TestData, 1)
	if err != nil {
		t.Fatal(err)
	}
	data, err = c.Get(TestKey)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(data, TestData) != 0 {
		t.Fatal(data)
	}
	err = c.Delete(TestKey)
	if err != nil {
		t.Fatal(err)
	}
	data, err = c.Get(TestKey)
	if err != herbdata.ErrNotFound {
		t.Fatal(err)
	}

	err = c.SetWithTTL(TestKey, TestData, 1)
	if err != nil {
		t.Fatal(err)
	}
	data, err = c.Get(TestKey)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(data, TestData) != 0 {
		t.Fatal(data)
	}
	time.Sleep(2 * time.Second)
	data, err = c.Get(TestKey)
	if err != herbdata.ErrNotFound {
		t.Fatal(err)
	}
	err = c.SetWithTTL(TestKey, TestData, 3600)
	if err != nil {
		t.Fatal(err)
	}
	data, err = c.Get(TestKey)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(data, TestData) != 0 {
		t.Fatal(data)
	}
	err = c.Revoke()
	if err != nil {
		t.Fatal(err)
	}
	data, err = c.Get(TestKey)
	if err != herbdata.ErrNotFound {
		t.Fatal(err)
	}
	err = c.SetWithTTL(TestKey, TestData, 3600)
	if err != nil {
		t.Fatal(err)
	}
	data, err = c.Get(TestKey)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(data, TestData) != 0 {
		t.Fatal(data)
	}
	err = c.Delete(TestKey)
	if err != nil {
		t.Fatal(err)
	}
	data, err = c.Get(TestKey)
	if err != herbdata.ErrNotFound {
		t.Fatal(err)
	}

}

func TestIrrevocableCache(t *testing.T) {
	var err error
	var data []byte
	c := newTestCache().WithIrrevocable(true)
	if !c.Irrevocable() {
		t.Fatal(c.Irrevocable())
	}
	err = c.SetWithTTL(TestKey, TestData, 1)
	if err != nil {
		t.Fatal(err)
	}
	data, err = c.Get(TestKey)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(data, TestData) != 0 {
		t.Fatal(data)
	}
	err = c.Delete(TestKey)
	if err != nil {
		t.Fatal(err)
	}
	data, err = c.Get(TestKey)
	if err != herbdata.ErrNotFound {
		t.Fatal(err)
	}

	err = c.SetWithTTL(TestKey, TestData, 1)
	if err != nil {
		t.Fatal(err)
	}
	data, err = c.Get(TestKey)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(data, TestData) != 0 {
		t.Fatal(data)
	}
	time.Sleep(2 * time.Second)
	data, err = c.Get(TestKey)
	if err != herbdata.ErrNotFound {
		t.Fatal(err)
	}
	err = c.SetWithTTL(TestKey, TestData, 3600)
	if err != nil {
		t.Fatal(err)
	}
	data, err = c.Get(TestKey)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(data, TestData) != 0 {
		t.Fatal(data)
	}
	err = c.Revoke()
	if err != ErrCacheIrrevocable {
		t.Fatal(err)
	}
}

func TestSet(t *testing.T) {
	var err error
	var data []byte
	c := newTestCache().WithTTL(1)
	err = c.Set(TestKey, TestData)
	if err != nil {
		t.Fatal(err)
	}
	data, err = c.Get(TestKey)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(data, TestData) != 0 {
		t.Fatal(data)
	}
	time.Sleep(2 * time.Second)
	data, err = c.Get(TestKey)
	if err != herbdata.ErrNotFound {
		t.Fatal(err)
	}
}

func TestLocker(t *testing.T) {
	var result = ""
	c := newTestCache()
	c2 := c.Child([]byte("c2"))
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		c.Lock(TestKey)
		defer c.Unlock(TestKey)
		time.Sleep(5 * time.Millisecond)
		result = result + "1"
		wg.Done()
	}()
	go func() {
		time.Sleep(2 * time.Millisecond)
		c2.RLock(TestKey)
		c2.RUnlock(TestKey)
		result = result + "2"
		wg.Done()
	}()
	go func() {
		time.Sleep(1 * time.Millisecond)
		c2.RLock(TestKey2)
		c2.RUnlock(TestKey2)
		result = result + "3"
		wg.Done()
	}()
	wg.Wait()
	if result != "312" {
		t.Fatal(result)
	}
}
