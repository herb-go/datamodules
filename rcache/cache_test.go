package rcache

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

var _ herbdata.RevocableNamespacedCache = New()

func newTestCache() *Cache {
	c := New()
	config := &Config{
		Store: &kvdb.Config{
			Driver: "freecache",
			Config: func(v interface{}) error {
				return json.Unmarshal([]byte(`{"Size":50000}`), v)
			},
		},
	}

	err := config.ApplyTo(c)
	if err != nil {
		panic(err)
	}
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
	var namespace = []byte("namespace")
	c := newTestCache()
	defer c.Stop()
	if c.Irrevocable() {
		t.Fatal(c.Irrevocable())
	}
	err = c.SetWithTTLNamespaced(namespace, TestKey, TestData, 1)
	if err != nil {
		t.Fatal(err)
	}
	data, err = c.GetNamespaced(namespace, TestKey)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(data, TestData) != 0 {
		t.Fatal(data)
	}
	err = c.DeleteNamespaced(namespace, TestKey)
	if err != nil {
		t.Fatal(err)
	}
	data, err = c.GetNamespaced(namespace, TestKey)
	if err != herbdata.ErrNotFound {
		t.Fatal(err)
	}

	err = c.SetWithTTLNamespaced(namespace, TestKey, TestData, 1)
	if err != nil {
		t.Fatal(err)
	}
	data, err = c.GetNamespaced(namespace, TestKey)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(data, TestData) != 0 {
		t.Fatal(data)
	}
	time.Sleep(2 * time.Second)
	data, err = c.GetNamespaced(namespace, TestKey)
	if err != herbdata.ErrNotFound {
		t.Fatal(err)
	}
	err = c.SetWithTTLNamespaced(namespace, TestKey, TestData, 3600)
	if err != nil {
		t.Fatal(err)
	}
	data, err = c.GetNamespaced(namespace, TestKey)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(data, TestData) != 0 {
		t.Fatal(data)
	}
	err = c.RevokeNamespaced(namespace)
	if err != nil {
		t.Fatal(err)
	}
	data, err = c.GetNamespaced(namespace, TestKey)
	if err != herbdata.ErrNotFound {
		t.Fatal(err)
	}
	err = c.SetWithTTLNamespaced(namespace, TestKey, TestData, 3600)
	if err != nil {
		t.Fatal(err)
	}
	data, err = c.GetNamespaced(namespace, TestKey)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(data, TestData) != 0 {
		t.Fatal(data)
	}
	err = c.DeleteNamespaced(namespace, TestKey)
	if err != nil {
		t.Fatal(err)
	}
	data, err = c.GetNamespaced(namespace, TestKey)
	if err != herbdata.ErrNotFound {
		t.Fatal(err)
	}

}

func TestIrrevocableCache(t *testing.T) {
	var err error
	var data []byte
	var namespace = []byte("namespace")
	c := newTestCache().WithIrrevocable(true)
	defer c.Stop()
	if !c.Irrevocable() {
		t.Fatal(c.Irrevocable())
	}
	err = c.SetWithTTLNamespaced(namespace, TestKey, TestData, 1)
	if err != nil {
		t.Fatal(err)
	}
	data, err = c.GetNamespaced(namespace, TestKey)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(data, TestData) != 0 {
		t.Fatal(data)
	}
	err = c.DeleteNamespaced(namespace, TestKey)
	if err != nil {
		t.Fatal(err)
	}
	data, err = c.GetNamespaced(namespace, TestKey)
	if err != herbdata.ErrNotFound {
		t.Fatal(err)
	}

	err = c.SetWithTTLNamespaced(namespace, TestKey, TestData, 1)
	if err != nil {
		t.Fatal(err)
	}
	data, err = c.GetNamespaced(namespace, TestKey)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(data, TestData) != 0 {
		t.Fatal(data)
	}
	time.Sleep(2 * time.Second)
	data, err = c.GetNamespaced(namespace, TestKey)
	if err != herbdata.ErrNotFound {
		t.Fatal(err)
	}
	err = c.SetWithTTLNamespaced(namespace, TestKey, TestData, 3600)
	if err != nil {
		t.Fatal(err)
	}
	data, err = c.GetNamespaced(namespace, TestKey)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(data, TestData) != 0 {
		t.Fatal(data)
	}
	err = c.RevokeNamespaced(namespace)
	if err != herbdata.ErrIrrevocable {
		t.Fatal(err)
	}
}

func TestCacheOperations(t *testing.T) {
	var cc *Cache
	c := newTestCache()
	e := c.engine
	defer e.Stop()
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

	cc = c.WithIrrevocable(true)
	if cc.irrevocable != true || c.irrevocable != false {
		t.Fatal(cc, c)
	}
	cc = c.WithEngine(nil)
	if cc.engine != nil || c.engine == nil {
		t.Fatal(cc, c)
	}
	cc = c.WithEngine(nil).WithIrrevocable(true)
	c.CopyFrom(cc)
	if c.irrevocable != true ||
		c.engine != nil {
		t.Fatal(c)
	}
}

func TestNestedCache(t *testing.T) {
	var err error
	var data []byte
	c := newTestCache()
	defer c.Stop()

	err = c.SetWithTTLNamespaced([]byte("namespace"), []byte("testkey2"), []byte("testvalue2"), 3600)
	if err != nil {
		panic(err)
	}
	data, err = c.GetNamespaced([]byte("namespace"), []byte("testkey2"))
	if err != nil || string(data) != "testvalue2" {
		t.Fatal()
	}
	err = c.SetWithTTLNamespaced([]byte("namespace2"), []byte("testkey2"), []byte("testvalue2n2"), 3600)
	if err != nil {
		panic(err)
	}
	data, err = c.GetNamespaced([]byte("namespace2"), []byte("testkey2"))
	if err != nil || string(data) != "testvalue2n2" {
		t.Fatal()
	}
	data, err = c.GetNamespaced([]byte("namespace"), []byte("testkey2"))
	if err != nil || string(data) != "testvalue2" {
		t.Fatal()
	}
}
