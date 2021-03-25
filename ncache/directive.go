package ncache

type Directive interface {
	Execute(*Cache) error
}

type DirectiveFunc func(*Cache) error

func (f DirectiveFunc) Execute(c *Cache) error {
	return f(c)
}

type Flushable bool

func (f Flushable) Execute(c *Cache) error {
	SetCache(c, c.VaryFlushable(bool(f)))
	return nil
}

func Child(path []byte) Directive {
	return DirectiveFunc(func(c *Cache) error {
		SetCache(c, c.Child(path))
		return nil
	})
}
func Migrate(namespace []byte) Directive {
	return DirectiveFunc(func(c *Cache) error {
		SetCache(c, c.Migrate(namespace))
		return nil
	})
}
func Prefix(prefix []byte) Directive {
	return DirectiveFunc(func(c *Cache) error {
		SetCache(c, c.VaryPrefix(prefix))
		return nil
	})
}
func SetCacheGroup(c *Cache, group []byte) {
	c.group = group
}
func SetCacheNamespace(c *Cache, namespace []byte) {
	c.namespace = namespace
}
func SetCachePath(c *Cache, path *Path) {
	c.path = path
}

func SetCacheFlushable(c *Cache, flushable bool) {
	c.flushable = flushable
}

func SetCacheStorage(c *Cache, storage *Storage) {
	c.storage = storage
}
func SetCachePromises(c *Cache, promises ...Directive) {
	c.promises = promises
}
func SetCache(c *Cache, src *Cache) {
	*c = *src
}
