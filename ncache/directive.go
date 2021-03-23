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

func Child(path ...[]byte) Directive {
	return DirectiveFunc(func(c *Cache) error {
		SetCache(c, c.Child(path...))
		return nil
	})
}

func Namespace(ns ...[]byte) Directive {
	return DirectiveFunc(func(c *Cache) error {
		SetCache(c, c.VaryNamesapce(ns...))
		return nil
	})
}

func Suffix(suffixs ...[]byte) Directive {
	return DirectiveFunc(func(c *Cache) error {
		SetCache(c, c.VarySuffix(suffixs...))
		return nil
	})
}

func SetCacheNamespaceTree(c *Cache, nst [][]byte) {
	c.namespaceTree = nst
}

func SetCacheRevocable(c *Cache, revocable bool) {
	c.revocable = revocable
}

func SetCacheStorage(c *Cache, storage *Storage) {
	c.storage = storage
}
func SetCacheTodos(c *Cache, todos ...Directive) {
	c.todos = todos
}
func SetCache(c *Cache, src *Cache) {
	*c = *src
}
