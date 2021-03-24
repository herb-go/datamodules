package ncache

type Directive interface {
	Execute(*Cache) error
}

type DirectiveFunc func(*Cache) error

func (f DirectiveFunc) Execute(c *Cache) error {
	return f(c)
}

type Revocable bool

func (r Revocable) Execute(c *Cache) error {
	SetCache(c, c.VaryRevocable(bool(r)))
	return nil
}

func Child(path []byte) Directive {
	return DirectiveFunc(func(c *Cache) error {
		SetCache(c, c.Child(path))
		return nil
	})
}
func Namespace(namespace []byte) Directive {
	return DirectiveFunc(func(c *Cache) error {
		SetCache(c, c.VaryNamesapce(namespace))
		return nil
	})
}
func Prefix(prefix []byte) Directive {
	return DirectiveFunc(func(c *Cache) error {
		SetCache(c, c.VaryPrefix(prefix))
		return nil
	})
}
func SetCachePrefix(c *Cache, prefix []byte) {
	c.prefix = prefix
}
func SetCacheNamespace(c *Cache, namespace []byte) {
	c.namespace = namespace
}
func SetCachePath(c *Cache, path *Path) {
	c.path = path
}

func SetCacheRevocable(c *Cache, revocable bool) {
	c.revocable = revocable
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
