package herbcache

type Config struct {
	Base *Config
	Directive
}

func (c *Config) Extend(d ...Directive) *Config {
	var extended *Config
	extended = c
	for _, v := range d {
		extended = &Config{
			Base:      extended,
			Directive: v,
		}
	}
	return extended
}

func (c *Config) ApplyTo(cache *Cache) error {
	if c == nil {
		return nil
	}
	err := c.Base.ApplyTo(cache)
	if err != nil {
		return err
	}
	return c.Directive.Execute(cache)
}
func NewConfig(d ...Directive) *Config {
	var c *Config
	return c.Extend(d...)
}

func LazyCache(d ...Directive) *Cache {
	c := New()
	c.config = NewConfig(d...)
	return c
}
