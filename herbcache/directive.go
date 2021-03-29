package herbcache

type Directive interface {
	Execute(*Cache) error
}

type DirectiveFunc func(*Cache) error

func (f DirectiveFunc) Execute(c *Cache) error {
	return f(c)
}

type Flushable bool

func (f Flushable) Execute(c *Cache) error {
	SetCache(c, c.OverrideFlushable(bool(f)))
	return nil
}

func SubCache(name []byte) Directive {
	return DirectiveFunc(func(c *Cache) error {
		SetCache(c, c.SubCache(name))
		return nil
	})
}
func Migrate(namespace []byte) Directive {
	return DirectiveFunc(func(c *Cache) error {
		SetCache(c, c.Migrate(namespace))
		return nil
	})
}
func Group(group []byte) Directive {
	return DirectiveFunc(func(c *Cache) error {
		SetCache(c, c.OverrideGroup(group))
		return nil
	})
}
