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
	return New().VaryRevocable(true).VaryStorage(s)
}

var TestKey = []byte("testkey")
var TestKey2 = []byte("testkey2")

var TestData = []byte("testdata")

func TestCache(t *testing.T) {
	var err error
	var data []byte
	var namespace = []byte("namespace")
	c := newTestCache().VarySuffix(namespace)
	c.storage.VersionTTL = 0
	defer c.Storage().Stop()
	if !c.Revocable() {
		t.Fatal(c.Revocable())
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

func TestCachedVersionCache(t *testing.T) {
	var err error
	var data []byte
	var namespace = []byte("namespace")
	c := newTestCache().VarySuffix(namespace)
	defer c.Storage().Stop()
	if !c.Revocable() {
		t.Fatal(c.Revocable())
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
	c := newTestCache().VaryRevocable(false).VarySuffix([]byte("namespace"))
	defer c.Storage().Stop()
	if c.Revocable() {
		t.Fatal(c.Revocable())
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

	cc = c.VaryRevocable(false)
	if cc.revocable != false || c.revocable != true {
		t.Fatal(cc, c)
	}
	cc = c.VaryStorage(nil)
	if cc.storage != nil || c.storage == nil {
		t.Fatal(cc, c)
	}
	cc = c.VarySuffix([]byte("suffix"))
	if cc.Equal(c) {
		t.Fatal(cc, c)
	}
	cc = c.Child([]byte("child"))
	if cc.Equal(c) {
		t.Fatal(cc, c)
	}
	cc = c.VaryStorage(nil).VaryRevocable(false)
	c.CopyFrom(cc)
	if c.revocable != false ||
		c.storage != nil {
		t.Fatal(c)
	}
}

func TestNoVersionStoreCache(t *testing.T) {
	var err error
	c := newTestCache()
	c.storage.VersionStore = nil
	defer c.Storage().Stop()
	err = c.Revoke()
	if err != ErrNoVersionStore {
		t.Fatal()
	}
}
func TestNamespace(t *testing.T) {
	var err error
	var data []byte
	c := newTestCache()
	defer c.Storage().Stop()
	var testkey = []byte("testkey")
	ctest1 := c.VarySuffix([]byte("test1"))
	ctest2 := ctest1.VarySuffix([]byte("test2"))
	ctest2_test3 := ctest2.Child([]byte("test3"))
	ctest2_test3_test4 := ctest2_test3.Child([]byte("test4"))
	err = ctest2_test3_test4.SetWithTTL(testkey, []byte("test4"), 3600)
	if err != nil {
		panic(err)
	}
	err = ctest2_test3.SetWithTTL(testkey, []byte("test3"), 3600)
	if err != nil {
		panic(err)
	}
	err = ctest2.SetWithTTL(testkey, []byte("test2"), 3600)
	if err != nil {
		panic(err)
	}
	err = ctest1.SetWithTTL(testkey, []byte("test1"), 3600)
	if err != nil {
		panic(err)
	}
	data, err = ctest2_test3_test4.Get(testkey)
	if !bytes.Equal(data, []byte("test4")) || err != nil {
		t.Fatal(data, err)
	}
	data, err = ctest2.Child([]byte("test3"), []byte("test4")).Get(testkey)
	if !bytes.Equal(data, []byte("test4")) || err != nil {
		t.Fatal(data, err)
	}
	data, err = ctest2.Child([]byte("test3test4")).Get(testkey)
	if err != herbdata.ErrNotFound {
		t.Fatal(data, err)
	}
	data, err = ctest2_test3.Get(testkey)
	if !bytes.Equal(data, []byte("test3")) || err != nil {
		t.Fatal(data, err)
	}
	data, err = ctest2.Get(testkey)
	if !bytes.Equal(data, []byte("test2")) || err != nil {
		t.Fatal(data, err)
	}
	data, err = c.VarySuffix([]byte("test1")).VarySuffix([]byte("test2")).Get(testkey)
	if !bytes.Equal(data, []byte("test2")) || err != nil {
		t.Fatal(data, err)
	}
	data, err = c.VarySuffix([]byte("test1test2")).Get(testkey)
	if err != herbdata.ErrNotFound {
		t.Fatal(data, err)
	}

	data, err = ctest1.Get(testkey)
	if !bytes.Equal(data, []byte("test1")) || err != nil {
		t.Fatal(data, err)
	}
	err = ctest1.Revoke()
	if err != nil {
		panic(err)
	}
	data, err = ctest1.Get(testkey)
	if err != herbdata.ErrNotFound {
		t.Fatal(data, err)
	}
	data, err = ctest2.Get(testkey)
	if !bytes.Equal(data, []byte("test2")) || err != nil {
		t.Fatal(data, err)
	}
	data, err = ctest2_test3_test4.Get(testkey)
	if !bytes.Equal(data, []byte("test4")) || err != nil {
		t.Fatal(data, err)
	}
	data, err = ctest2.Child([]byte("test3"), []byte("test4")).Get(testkey)
	if !bytes.Equal(data, []byte("test4")) || err != nil {
		t.Fatal(data, err)
	}
	err = ctest2.Revoke()
	if err != nil {
		panic(err)
	}
	data, err = ctest2.Get(testkey)
	if err != herbdata.ErrNotFound {
		t.Fatal(data, err)
	}
	data, err = ctest2_test3_test4.Get(testkey)
	if err != herbdata.ErrNotFound {
		t.Fatal(data, err)
	}
	data, err = ctest2.Child([]byte("test3"), []byte("test4")).Get(testkey)
	if err != herbdata.ErrNotFound {
		t.Fatal(data, err)
	}
	ctestnamespace := c.VaryNamesapce([]byte("ns1"), []byte("ns2"), []byte("ns3"))
	ctestsuffix := c.VaryNamesapce([]byte("ns1")).VarySuffix([]byte("ns2"), []byte("ns3"))
	ctestsuffix2 := c.VaryNamesapce([]byte("ns1")).VarySuffix([]byte("ns2")).VarySuffix([]byte("ns3"))
	ctestsuffix3 := c.VarySuffix([]byte("ns1"), []byte("ns2"), []byte("ns3"))
	ctestnamespacefail := c.VaryNamesapce([]byte("ns1")).VaryNamesapce([]byte("ns2"), []byte("ns3"))
	for _, v := range []*Cache{ctestnamespace, ctestsuffix, ctestsuffix2, ctestsuffix3, ctestnamespacefail} {
		data, err = v.Get(testkey)
		if err != herbdata.ErrNotFound {
			t.Fatal(data, err)
		}
	}
	err = ctestnamespace.SetWithTTL(testkey, []byte("testdata"), 3600)
	if err != nil {
		panic(err)
	}
	for _, v := range []*Cache{ctestnamespace, ctestsuffix, ctestsuffix2, ctestsuffix3} {
		data, err = v.Get(testkey)
		if !bytes.Equal(data, []byte("testdata")) || err != nil {
			t.Fatal(data, err)
		}
	}
	for _, v := range []*Cache{ctestnamespacefail} {
		data, err = v.Get(testkey)
		if err != herbdata.ErrNotFound {
			t.Fatal(data, err)
		}
	}
}
