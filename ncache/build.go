package ncache

type Builder interface {
	BuildCache(*NestedCache) error
}

type BuilderFunc func(*NestedCache) error

func (f BuilderFunc) BuildCache(c *NestedCache) error {
	return f(c)
}

type NestRevocable bool

func (r NestRevocable) BuildCache(c *NestedCache) error {
	c.Cache = c.Cache.WithRevocable(bool(r))
	return nil
}

func NestChild(path ...[]byte) Builder {
	return BuilderFunc(func(c *NestedCache) error {
		c.Cache = c.Cache.Child(path...)
		return nil
	})
}

func NestNamespace(ns ...[]byte) Builder {
	return BuilderFunc(func(c *NestedCache) error {
		c.Cache = c.Cache.WithNamesapce(ns...)
		return nil
	})
}

func NestSuffix(suffixs ...[]byte) Builder {
	return BuilderFunc(func(c *NestedCache) error {
		c.Cache = c.Cache.WithSuffix(suffixs...)
		return nil
	})
}

func NestEngine(engine *Engine) Builder {
	return BuilderFunc(func(c *NestedCache) error {
		c.Cache = c.Cache.WithEngine(engine)
		return nil
	})
}
