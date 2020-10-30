package rcache

import (
	"bytes"
	"time"

	"github.com/herb-go/herbdata/dataencoding"

	"github.com/herb-go/herbdata/datautil"
	"github.com/herb-go/herbdata/kvdb"
)

type CachePathPrefix []byte

func (p CachePathPrefix) Join(pathlist ...[]byte) []byte {
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
	path        []byte
	ttl         time.Duration
	irrevocable bool
	encoding    *dataencoding.Encoding
	engine      *Engine
}

func (c *Cache) Irrevocable() bool {
	return c.irrevocable
}
func (c *Cache) Revoke() error {
	v, err := c.engine.VersionGenerator()
	if err != nil {
		return err
	}
	return c.engine.VersionStore.Set(CachePathPrefixVersion.Join(c.path), []byte(v))
}
func (c *Cache) Get(key []byte) ([]byte, error) {
	var data []byte
	var version []byte
	var revocable bool
	var err error
	var e *enity
	if len(key) == 0 {
		return nil, kvdb.ErrInvalidateKey
	}
	if !c.irrevocable {
		revocable = true
		version, err = c.engine.VersionStore.Get(CachePathPrefixVersion.Join(c.path))
		if err != nil {
			return nil, err
		}
	}
	data, err = c.engine.Store.Get(CachePathPrefixValue.Join(c.path, key))
	if err != nil {
		return nil, err
	}
	e, err = loadEnity(data, revocable, version)
	if err != nil {
		if err == ErrEnityTypecodeNotMatch || err == ErrEnityVersionNotMatch {
			return nil, kvdb.ErrKeyNotFound
		}
		return nil, err
	}
	return e.data, nil

}
func (c *Cache) Set(key []byte, data []byte) error {
	return c.SetWithTTL(key, data, c.ttl)
}
func (c *Cache) SetWithTTL(key []byte, data []byte, ttl time.Duration) error {
	var version []byte
	var revocable bool
	var err error
	var e *enity
	if !c.irrevocable {
		revocable = true
		version, err = c.engine.VersionStore.Get(CachePathPrefixVersion.Join(c.path))
		if err != nil {
			return err
		}
	}
	e = createEnity(revocable, version, data)
	buf := bytes.NewBuffer(nil)
	err = e.SaveTo(buf)
	if err != nil {
		return err
	}
	return c.engine.Store.SetWithTTL(CachePathPrefixValue.Join(c.path, key), buf.Bytes(), ttl)
}

func (c *Cache) Del(key []byte) error {
	return c.engine.Store.Del(CachePathPrefixValue.Join(c.path, key))
}
func (c *Cache) Clone() *Cache {
	return &Cache{
		path:        c.path,
		ttl:         c.ttl,
		irrevocable: c.irrevocable,
		engine:      c.engine,
		encoding:    c.encoding,
	}
}

func (c *Cache) Child(path []byte) *Cache {
	cc := c.Clone()
	cc.path = datautil.Join(c.path, path)
	return cc
}

func (c *Cache) WithIrrevocable(irrevocable bool) *Cache {
	cc := c.Clone()
	cc.irrevocable = irrevocable
	return cc
}
func (c *Cache) WithPath(path []byte) *Cache {
	cc := c.Clone()
	cc.path = path
	return cc
}

func (c *Cache) WithEngine(engine *Engine) *Cache {
	cc := c.Clone()
	cc.engine = engine
	return cc
}

func (c *Cache) WithTTL(ttl time.Duration) *Cache {
	cc := c.Clone()
	cc.ttl = ttl
	return cc
}
func (c *Cache) WithEncoding(e *dataencoding.Encoding) *Cache {
	cc := c.Clone()
	cc.encoding = e
	return cc
}
func (c *Cache) CopyFrom(src *Cache) {
	c.path = src.path
	c.ttl = src.ttl
	c.irrevocable = src.irrevocable
	c.engine = src.engine
}

func (c *Cache) RLock(key []byte) {
	c.engine.lockerMap.RLock(string(key))
}
func (c *Cache) RUnlock(key []byte) {
	c.engine.lockerMap.RUnlock(string(key))
}
func (c *Cache) Lock(key []byte) {
	c.engine.lockerMap.Lock(string(key))
}
func (c *Cache) Unlock(key []byte) {
	c.engine.lockerMap.Unlock(string(key))
}

func New() *Cache {
	return &Cache{}
}
