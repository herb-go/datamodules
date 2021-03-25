package herbcache

type Directive interface {
	Execute(*Cache) error
}
