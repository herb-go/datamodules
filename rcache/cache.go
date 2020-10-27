package rcache

import (
	"bytes"
	"time"

	"github.com/herb-go/herbdata/datautil"
	"github.com/herb-go/herbdata/kvdb"
)

type Cache struct {
	Path        []byte
	irrevocable bool
	engine      *Engine
}

func (c *Cache) Irrevocable() bool {
	return c.irrevocable
}
func (c *Cache) Revoke() error {
	v, err := c.engine.versionGenerator()
	if err != nil {
		return err
	}
	return c.engine.versionstore.SetVersion(c.Path, []byte(v))
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
		version, err = c.engine.versionstore.GetVersion(c.Path)
		if err != nil {
			return nil, err
		}
	}
	data, err = c.engine.driver.Get(datautil.Join(nil, c.Path, key))
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

func (c *Cache) SetWithTTL(key []byte, data []byte, ttl time.Duration) error {
	var version []byte
	var revocable bool
	var err error
	var e *enity
	if !c.irrevocable {
		revocable = true
		version, err = c.engine.versionstore.GetVersion(c.Path)
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
	return c.engine.driver.SetWithTTL(datautil.Join(nil, c.Path, key), buf.Bytes(), ttl)
}

func (c *Cache) Del(key []byte) error {
	return c.engine.driver.Del(datautil.Join(c.Path, key))
}
func (c *Cache) NewCache(path []byte, irrevocable bool) *Cache {
	return &Cache{
		Path:        datautil.Join(c.Path, path),
		irrevocable: irrevocable,
		engine:      c.engine,
	}
}
