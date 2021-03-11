package rcache

import (
	"bytes"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/herb-go/herbdata"

	_ "github.com/herb-go/herbdata-drivers/kvdb-drivers/freecachedb"
	"github.com/herb-go/herbdata/dataencoding/jsonencoding"
	"github.com/herb-go/herbdata/kvdb"
	_ "github.com/herb-go/herbdata/kvdb/commonkvdb"
)

var _ herbdata.Cache = New()

func newTestCache() *Cache {
	c := New()
	config := &Config{
		Store: &kvdb.Config{
			Driver: "freecache",
			Config: func(v interface{}) error {
				return json.Unmarshal([]byte(`{"Size":50000}`), v)
			},
		},
		VersionStore: &kvdb.Config{
			Driver: "inmemory",
		},
	}

	err := config.ApplyTo(c)
	if err != nil {
		panic(err)
	}
	c.encoding = jsonencoding.Encoding
	err = c.Start()
	if err != nil {
		panic(err)
	}
	return c
}

var TestKey = []byte("testkey")
var TestKey2 = []byte("testkey2")

var TestData = []byte("testdata")

func TestCache(t *testing.T) {
	var err error
	var data []byte
	c := newTestCache()
	defer c.Stop()
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
	defer c.Stop()
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
	if err != herbdata.ErrIrrevocable {
		t.Fatal(err)
	}
}

func TestSet(t *testing.T) {
	var err error
	var data []byte
	c := newTestCache().WithTTL(1)
	defer c.Stop()
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
	defer c.Stop()
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

func TestCacheOperations(t *testing.T) {
	var cc *Cache
	c := newTestCache()
	e := c.engine
	defer e.Stop()
	c.ttl = 200
	c.path = []byte("path")
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
	cc.ttl = 100
	if c.ttl != 200 {
		t.Fatal(c.ttl)
	}
	cc = c.WithPath([]byte("newpath"))
	if bytes.Compare(cc.path, []byte("newpath")) != 0 {
		t.Fatal(cc)
	}
	if bytes.Compare(c.path, []byte("path")) != 0 {
		t.Fatal(c)
	}
	cc = c.Child([]byte("child"))
	if bytes.Compare(cc.path, []byte("path")) == 0 || bytes.Compare(cc.path[:len([]byte("path"))], []byte("path")) != 0 {
		t.Fatal(cc)
	}
	if bytes.Compare(c.path, []byte("path")) != 0 {
		t.Fatal(c)
	}
	cc = c.WithTTL(500)
	if cc.ttl != 500 || c.ttl != 200 {
		t.Fatal(cc, c)
	}
	cc = c.WithIrrevocable(true)
	if cc.irrevocable != true || c.irrevocable != false {
		t.Fatal(cc, c)
	}
	cc = c.WithEncoding(nil)
	if cc.encoding != nil || c.encoding == nil {
		t.Fatal(cc, c)
	}
	cc = c.WithEngine(nil)
	if cc.engine != nil || c.engine == nil {
		t.Fatal(cc, c)
	}
	cc = c.WithPath([]byte("newpath")).WithTTL(500).WithIrrevocable(true).WithEncoding(nil).WithEngine(nil)
	if bytes.Compare(c.path, []byte("path")) != 0 ||
		c.ttl != 200 ||
		c.irrevocable != false ||
		c.encoding == nil ||
		c.engine == nil {
		t.Fatal(c)
	}
	c.CopyFrom(cc)
	if bytes.Compare(c.path, []byte("newpath")) != 0 ||
		c.ttl != 500 ||
		c.irrevocable == false ||
		c.encoding != nil ||
		c.engine != nil {
		t.Fatal(c)
	}
}

func TestNestedCache(t *testing.T) {
	var err error
	var data []byte
	c := newTestCache()
	defer c.Stop()
	data, err = c.Child([]byte("testpath"), []byte("testpath2")).Get([]byte("testkey"))
	if len(data) != 0 || err != herbdata.ErrNotFound {
		t.Fatal()
	}
	child := c.Child([]byte("testpath"))
	child2 := child.Child([]byte("testpath2"))
	err = child2.SetWithTTL([]byte("testkey"), []byte("testvalue"), 3600)
	data, err = c.Child([]byte("testpath"), []byte("testpath2")).Get([]byte("testkey"))
	if err != nil || string(data) != "testvalue" {
		t.Fatal()
	}
	nestedpath := [][]byte{
		[]byte("testpath"), []byte("testpath2"),
	}
	err = c.SetWithTTLNested([]byte("testkey2"), []byte("testvalue2"), 3600, nestedpath...)
	if err != nil {
		panic(err)
	}
	data, err = c.GetNested([]byte("testkey2"), nestedpath...)
	if err != nil || string(data) != "testvalue2" {
		t.Fatal()
	}
	err = c.DeleteNested([]byte("testkey2"), nestedpath...)
	if err != nil {
		panic(err)
	}
	data, err = c.GetNested([]byte("testkey2"), nestedpath...)
	if len(data) != 0 || err != herbdata.ErrNotFound {
		t.Fatal()
	}
	err = c.SetWithTTLNested([]byte("testkey2"), []byte("testvalue2"), 3600, nestedpath...)
	if err != nil {
		panic(err)
	}
	data, err = c.GetNested([]byte("testkey2"), nestedpath...)
	if err != nil || string(data) != "testvalue2" {
		t.Fatal()
	}
	err = c.RevokeNested(nestedpath...)
	if err != nil {
		panic(err)
	}
	data, err = c.GetNested([]byte("testkey2"), nestedpath...)
	if len(data) != 0 || err != herbdata.ErrNotFound {
		t.Fatal()
	}

}
