package cachecommand

type Command interface {
	Exec(ctx *Context) error
}

type Preset struct {
	prev    *Preset
	command Command
}
