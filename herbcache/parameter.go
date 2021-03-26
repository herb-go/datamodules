package herbcache

type Parameter interface {
	Namespace() []byte
	Group() []byte
	Position() *Position
	Flushable() bool
}
