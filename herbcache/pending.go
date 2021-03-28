package herbcache

type Pending struct {
	Base *Pending
	Directive
}

func (p *Pending) Extend(d ...Directive) *Pending {
	var extended *Pending
	extended = p
	for _, v := range d {
		extended = &Pending{
			Base:      extended,
			Directive: v,
		}
	}
	return extended
}

func (c *Pending) Resolve(cache *Cache) error {
	if c == nil {
		return nil
	}
	err := c.Base.Resolve(cache)
	if err != nil {
		return err
	}
	return c.Directive.Execute(cache)
}
func Pend(d ...Directive) *Pending {
	var c *Pending
	return c.Extend(d...)
}
