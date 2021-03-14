package rcache

import (
	"bytes"

	"github.com/herb-go/herbdata"
	"github.com/herb-go/herbdata/datautil"
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
	irrevocable bool
	engine      *Engine
}

func (c *Cache) Irrevocable() bool {
	return c.irrevocable
}
func (c *Cache) getVersion(namespace []byte) ([]byte, error) {
	version, err := c.engine.Store.Get(CachePathPrefixVersion.Join(namespace))
	if err != nil {
		if err == herbdata.ErrNotFound {
			return []byte{}, nil
		}
		return nil, err

	}
	return version, nil
}
func (c *Cache) setVersion(namespace []byte, version []byte) error {
	return c.engine.Store.SetWithTTL(CachePathPrefixVersion.Join(namespace), version, c.engine.VersionTTL)
}

func (c *Cache) RevokeNamespaced(namespace []byte) error {
	if c.irrevocable {
		return herbdata.ErrIrrevocable
	}
	v, err := c.engine.VersionGenerator()
	if err != nil {
		return err
	}
	return c.setVersion(namespace, []byte(v))
}
func (c *Cache) GetNamespaced(namespace []byte, key []byte) ([]byte, error) {
	var data []byte
	var version []byte
	var revocable bool
	var err error
	var e *enity
	if !c.irrevocable {
		revocable = true
		version, err = c.getVersion(namespace)
		if err != nil {
			return nil, err
		}
	}
	data, err = c.engine.Store.Get(CachePathPrefixValue.Join(namespace, key))
	if err != nil {
		return nil, err
	}
	e, err = loadEnity(data, revocable, version)
	if err != nil {
		if err == ErrEnityTypecodeNotMatch || err == ErrEnityVersionNotMatch {
			return nil, herbdata.ErrNotFound
		}
		return nil, err
	}
	return e.data, nil

}

func (c *Cache) SetWithTTLNamespaced(namespace []byte, key []byte, data []byte, ttl int64) error {
	var version []byte
	var revocable bool
	var err error
	var e *enity
	if !c.irrevocable {
		revocable = true
		version, err = c.getVersion(namespace)
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
	return c.engine.Store.SetWithTTL(CachePathPrefixValue.Join(namespace, key), buf.Bytes(), ttl)
}

func (c *Cache) DeleteNamespaced(namespace []byte, key []byte) error {
	return c.engine.Store.Delete(CachePathPrefixValue.Join(namespace, key))
}

func (c *Cache) Clone() *Cache {
	return &Cache{
		irrevocable: c.irrevocable,
		engine:      c.engine,
	}
}

func (c *Cache) WithIrrevocable(irrevocable bool) *Cache {
	cc := c.Clone()
	cc.irrevocable = irrevocable
	return cc
}

func (c *Cache) WithEngine(engine *Engine) *Cache {
	cc := c.Clone()
	cc.engine = engine
	return cc
}

func (c *Cache) CopyFrom(src *Cache) {
	c.irrevocable = src.irrevocable
	c.engine = src.engine
}
func (c *Cache) Equal(dst *Cache) bool {
	if dst == nil || c == nil {
		return dst == nil && c == nil
	}
	return c.irrevocable == dst.irrevocable &&
		c.engine == dst.engine
}
func (c *Cache) Start() error {
	return c.engine.Start()
}
func (c *Cache) Stop() error {
	return c.engine.Stop()
}
func New() *Cache {
	return &Cache{}
}
