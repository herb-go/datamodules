package herbcache

import "bytes"

type Cache struct {
	storage   Storage
	config    *Config
	namespace []byte
	group     []byte
	position  *Position
	flushable bool
}

func (c *Cache) Get(key []byte) ([]byte, error) {
	return c.storage.ExecuteGet(c, key)
}
func (c *Cache) SetWithTTL(key []byte, data []byte, ttl int64) error {
	return c.storage.ExecuteSetWithTTL(c, key, data, ttl)
}
func (c *Cache) Delete(key []byte) error {
	return c.storage.ExecuteDelete(c, key)
}
func (c *Cache) Flush() error {
	return c.storage.ExecuteFlush(c)
}

func (c *Cache) Clone() *Cache {
	return &Cache{
		storage:   c.storage,
		namespace: c.namespace,
		group:     c.group,
		position:  c.position,
		flushable: c.flushable,
		config:    c.config,
	}
}
func (c *Cache) Equal(dst *Cache) bool {
	if dst == nil || c == nil {
		return dst == c
	}
	if !c.position.Equal(dst.position) {
		return false
	}
	if !bytes.Equal(c.group, dst.group) {
		return false
	}
	if !bytes.Equal(c.namespace, dst.namespace) {
		return false
	}
	return c.flushable == dst.flushable &&
		c.storage == dst.storage
}

func (c *Cache) Migrate(namespace []byte) *Cache {
	cc := c.Clone()
	SetCacheNamespace(cc, namespace)
	SetCacheGroup(cc, nil)
	return cc
}

func (c *Cache) Namespace() []byte {
	return c.namespace
}

func (c *Cache) OverrideGroup(group []byte) *Cache {
	cc := c.Clone()
	SetCacheGroup(cc, group)
	return cc
}

func (c *Cache) Group() []byte {
	return c.group
}

func (c *Cache) Position() *Position {
	return c.position
}

func (c *Cache) OverrideFlushable(flushable bool) *Cache {
	cc := c.Clone()
	SetCacheFlushable(cc, flushable)
	return cc
}

func (c *Cache) Flushable() bool {
	return c.flushable
}

func (c *Cache) IsPreparing() bool {
	return c.config == nil
}

func (c *Cache) Ready() error {
	conf := c.config
	c.config = nil
	return conf.ApplyTo(c)
}
func New() *Cache {
	return &Cache{}
}
func Prepare(d ...Directive) *Cache {
	c := New()
	c.config = NewConfig(d...)
	return c
}

func SetCacheStorage(c *Cache, s Storage) {
	c.storage = s
}
func SetCacheNamespace(c *Cache, namespace []byte) {
	c.namespace = namespace
}
func SetCachePosition(c *Cache, position *Position) {
	c.position = position
}

func SetCacheGroup(c *Cache, group []byte) {
	c.group = group
}

func SetCacheFlushable(c *Cache, flushable bool) {
	c.flushable = flushable
}

func SetCache(c *Cache, dst *Cache) {
	*c = *dst
}
func SetCacheConfig(c *Cache, p *Config) {
	c.config = p
}
func Copy(src *Cache, dst *Cache) {
	SetCache(src, dst.Clone())
}
