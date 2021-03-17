package ncache

type NestedCache struct {
	*Cache
	Builders []Builder
}

func (c *NestedCache) WithBuilder(b ...Builder) *NestedCache {
	c.Builders = append(c.Builders, b...)
	return c
}
func (c *NestedCache) BuildCache(nested *NestedCache) error {
	for _, v := range c.Builders {
		err := v.BuildCache(nested)
		if err != nil {
			return err
		}
	}
	return nil
}
func (c *NestedCache) Start() error {
	for _, v := range c.Builders {
		err := v.BuildCache(c)
		if err != nil {
			return err
		}
	}
	return nil
}
func (c *NestedCache) Stop() error {
	return nil
}
func (c *NestedCache) NewNested(builder ...Builder) *NestedCache {
	return NewNestedCache(c).WithBuilder(builder...)
}

func NewNestedCache(b ...Builder) *NestedCache {
	c := &NestedCache{}
	return c.WithBuilder(b...)
}
