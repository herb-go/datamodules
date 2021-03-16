package versionedcache

import (
	"bytes"

	"github.com/herb-go/herbdata"
	"github.com/herb-go/herbdata/datautil"
)

type CachePathPrefix []byte

func (p CachePathPrefix) MustJoin(pathlist ...[]byte) []byte {
	buf := bytes.NewBuffer(nil)
	_, err := buf.Write(p)
	if err != nil {
		panic(err)
	}
	err = datautil.PackTo(buf, nil, pathlist...)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

var CachePathPrefixValue = CachePathPrefix([]byte{0})
var CachePathPrefixVersion = CachePathPrefix([]byte{1})

type Cache struct {
	revocable     bool
	NamespaceTree [][]byte
	engine        *Engine
}

func (c *Cache) Revocable() bool {
	return c.revocable
}
func (c *Cache) getCachedVersion(key []byte) ([]byte, error) {
	version, err := c.engine.VersionStore.Get(key)
	if err == nil {
		return version, nil
	}
	if err != herbdata.ErrNotFound {
		return nil, err
	}
	return c.engine.LoadRawVersion(key)

}
func (c *Cache) getVersion() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	var err error
	cacheable := c.engine.VersionTTL > 0 && c.engine.VersionStore != nil
	for k := range c.NamespaceTree {
		var version []byte
		key := CachePathPrefixVersion.MustJoin(c.NamespaceTree[0 : k+1]...)
		if cacheable {
			version, err = c.getCachedVersion(key)
		} else {
			version, err = c.engine.LoadRawVersion(key)
		}
		if err != nil {
			return nil, err
		}
		err = datautil.PackTo(buf, nil, version)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
func (c *Cache) setVersion(version []byte) error {
	cacheable := c.engine.VersionTTL > 0 && c.engine.VersionStore != nil
	key := CachePathPrefixVersion.MustJoin(c.NamespaceTree...)
	err := c.engine.VersionStore.Set(key, version)
	if err != nil {
		return err
	}
	if cacheable {
		return c.engine.Store.Delete(key)
	}
	return nil
}
func (c *Cache) mustGetNamespace() []byte {
	buf := bytes.NewBuffer(nil)
	err := datautil.PackTo(buf, nil, c.NamespaceTree...)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}
func (c *Cache) Revoke() error {
	if !c.revocable {
		return herbdata.ErrIrrevocable
	}
	if c.engine.VersionStore == nil {
		return ErrNoVersionStore
	}
	v, err := c.engine.VersionGenerator()
	if err != nil {
		return err
	}
	return c.setVersion([]byte(v))
}
func (c *Cache) Get(key []byte) ([]byte, error) {
	var data []byte
	var version []byte
	var err error
	var e *enity
	namespace := c.mustGetNamespace()
	if c.revocable {
		version, err = c.getVersion()
		if err != nil {
			return nil, err
		}
	}
	data, err = c.engine.Store.Get(CachePathPrefixValue.MustJoin(namespace, key))
	if err != nil {
		return nil, err
	}
	e, err = loadEnity(data, c.revocable, version)
	if err != nil {
		if err == ErrEnityTypecodeNotMatch || err == ErrEnityVersionNotMatch {
			return nil, herbdata.ErrNotFound
		}
		return nil, err
	}
	return e.data, nil

}

func (c *Cache) SetWithTTL(key []byte, data []byte, ttl int64) error {
	var version []byte
	var err error
	var e *enity
	namespace := c.mustGetNamespace()
	if c.revocable {
		version, err = c.getVersion()
		if err != nil {
			return err
		}
	}
	e = createEnity(c.revocable, version, data)
	buf := bytes.NewBuffer(nil)
	err = e.SaveTo(buf)
	if err != nil {
		return err
	}
	return c.engine.Store.SetWithTTL(CachePathPrefixValue.MustJoin(namespace, key), buf.Bytes(), ttl)
}

func (c *Cache) Delete(key []byte) error {
	namespace := c.mustGetNamespace()
	return c.engine.Store.Delete(CachePathPrefixValue.MustJoin(namespace, key))
}

func (c *Cache) Clone() *Cache {
	t := make([][]byte, len(c.NamespaceTree))
	for k := range t {
		t[k] = make([]byte, len(c.NamespaceTree[k]))
		copy(t[k], c.NamespaceTree[k])
	}
	return &Cache{
		revocable:     c.revocable,
		engine:        c.engine,
		NamespaceTree: t,
	}
}

func (c *Cache) WithRevocable(revocable bool) *Cache {
	cc := c.Clone()
	cc.revocable = revocable
	return cc
}
func (c *Cache) WithSuffix(suffix []byte) *Cache {
	index := len(c.NamespaceTree) - 1
	buf := bytes.NewBuffer(nil)
	_, err := buf.Write(c.NamespaceTree[index])
	if err != nil {
		panic(err)
	}
	err = datautil.PackTo(buf, nil, suffix)
	if err != nil {
		panic(err)
	}
	return c.WithNamesapce(buf.Bytes())
}
func (c *Cache) WithNamesapce(namespace []byte) *Cache {
	cc := c.Clone()
	index := len(cc.NamespaceTree) - 1
	cc.NamespaceTree[index] = namespace
	return cc
}
func (c *Cache) Child(name ...[]byte) *Cache {
	cc := c.Clone()
	cc.NamespaceTree = append(cc.NamespaceTree, name...)
	return cc
}
func (c *Cache) WithEngine(engine *Engine) *Cache {
	cc := c.Clone()
	cc.engine = engine
	return cc
}
func (c *Cache) Engine() *Engine {
	return c.engine
}
func (c *Cache) CopyFrom(src *Cache) {
	c.revocable = src.revocable
	c.engine = src.engine
}
func (c *Cache) Equal(dst *Cache) bool {
	if dst == nil || c == nil {
		return dst == nil && c == nil
	}
	if len(c.NamespaceTree) != len(dst.NamespaceTree) {
		return false
	}
	for k := range c.NamespaceTree {
		if bytes.Compare(c.NamespaceTree[k], dst.NamespaceTree[k]) != 0 {
			return false
		}
	}
	return c.revocable == dst.revocable &&
		c.engine == dst.engine
}
func (c *Cache) Start() error {
	return c.engine.Start()
}
func (c *Cache) Stop() error {
	return c.engine.Stop()
}
func New() *Cache {
	return &Cache{
		NamespaceTree: [][]byte{[]byte{}},
	}
}
