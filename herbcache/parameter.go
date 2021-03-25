package herbcache

type Parameter struct {
	namespace []byte
	group     []byte
	position  *Position
	flushable bool
}
