package herbcache

type Context interface {
	Namespace() []byte
	Group() []byte
	Position() *Position
	Flushable() bool
}
