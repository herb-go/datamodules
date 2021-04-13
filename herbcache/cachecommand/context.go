package cachecommand

type Context struct {
	preset Preset
}

func NewContext() *Context {
	return &Context{}
}
