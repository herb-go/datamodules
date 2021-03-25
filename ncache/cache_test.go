package ncache

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/herb-go/herbdata"

	_ "github.com/herb-go/herbdata-drivers/kvdb-drivers/freecachedb"
	"github.com/herb-go/herbdata/kvdb"
	_ "github.com/herb-go/herbdata/kvdb/commonkvdb"
)

var _ herbdata.NestableCache = New()

func newTestCache() *Cache {
	s := NewStorage()
	config := &StorageConfig{
		Cache: &kvdb.Config{
			Driver: "freecache",
			Config: func(v interface{}) error {
				return json.Unmarshal([]byte(`{"Size":50000}`), v)
			},
		},
		VersionTTL: 3600,
		VersionStore: &kvdb.Config{
			Driver: "inmemory",
		},
	}

	err := config.ApplyTo(s)
	if err != nil {
		panic(err)
	}
	err = s.Start()
	if err != nil {
		panic(err)
	}
	return New().VaryFlushable(true).VaryStorage(s)
}

var TestKey = []byte("testkey")
var TestKey2 = []byte("testkey2")

var TestData = []byte("testdata")

func TestCache(t *testing.T) {
	var err error
	var data []byte
	var namespace = []byte("namespace")
	c := newTestCache().VaryPrefix(namespace)
	c.storage.VersionTTL = 0
	defer c.Storage().Stop()
	if !c.Flushable() {
		t.Fatal(c.Flushable())
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
	err = c.Flush()
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

func TestCachedVersionCache(t *testing.T) {
	var err error
	var data []byte
	var namespace = []byte("namespace")
	c := newTestCache().VaryPrefix(namespace)
	defer c.Storage().Stop()
	if !c.Flushable() {
		t.Fatal(c.Flushable())
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
	err = c.Flush()
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
	c := newTestCache().VaryFlushable(false).VaryPrefix([]byte("namespace"))
	defer c.Storage().Stop()
	if c.Flushable() {
		t.Fatal(c.Flushable())
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
	err = c.Flush()
	if err != herbdata.ErrIrrevocable {
		t.Fatal(err)
	}
}

func TestCacheOperations(t *testing.T) {
	var cc *Cache
	c := newTestCache()
	s := c.storage
	defer s.Stop()
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

	cc = c.VaryFlushable(false)
	if cc.flushable != false || c.flushable != true {
		t.Fatal(cc, c)
	}
	cc = c.VaryStorage(nil)
	if cc.storage != nil || c.storage == nil {
		t.Fatal(cc, c)
	}
	cc = c.VaryPrefix([]byte("suffix"))
	if cc.Equal(c) {
		t.Fatal(cc, c)
	}
	cc = c.Child([]byte("child"))
	if cc.Equal(c) {
		t.Fatal(cc, c)
	}
	cc = c.VaryStorage(nil).VaryFlushable(false)
	c.CopyFrom(cc)
	if c.flushable != false ||
		c.storage != nil {
		t.Fatal(c)
	}
}

func TestNoVersionStoreCache(t *testing.T) {
	var err error
	c := newTestCache()
	c.storage.VersionStore = nil
	defer c.Storage().Stop()
	err = c.Flush()
	if err != ErrNoVersionStore {
		t.Fatal()
	}
}

func shouldEqual(t *testing.T, key []byte, target []byte, e error, caches ...*Cache) {
	for _, v := range caches {
		data, err := v.Get(key)
		if !bytes.Equal(data, target) || e != err {
			t.Fatal(data, err)
		}
	}
}
func shouldNotEqual(t *testing.T, key []byte, target []byte, e error, caches ...*Cache) {
	for _, v := range caches {
		data, err := v.Get(key)
		if bytes.Equal(data, target) || e == err {
			t.Fatal(data, err)
		}
	}
}
func must(err error) {
	if err != nil {
		panic(err)
	}
}
func TestNamespace(t *testing.T) {
	c := newTestCache()
	defer c.Storage().Stop()
	var testkey = []byte("testkey")
	var testdata = []byte("testdata")
	cns1 := c.Migrate([]byte("ns1"))
	cns2 := c.Migrate([]byte("ns2"))
	cns1sub1 := cns1.Child([]byte("sub1"))
	cns2sub1 := cns1.Child([]byte("sub1"))
	must(cns1.SetWithTTL(testkey, testdata, 3600))

	shouldEqual(t, testkey, testdata, nil, cns1)
	shouldNotEqual(t, testkey, testdata, nil, cns2, cns1sub1, cns2sub1)
	must(cns1.Delete(testkey))
	must(cns1sub1.SetWithTTL(testkey, testdata, 3600))
	shouldEqual(t, testkey, testdata, nil, cns1sub1)
	must(cns1.Flush())
	shouldNotEqual(t, testkey, testdata, nil, cns1sub1)

}
