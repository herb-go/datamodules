package cachecommand

type Preset struct {
	prev    *Preset
	command Command
}

func (p *Preset) concatCommand(cmd Command) *Preset {
	return &Preset{
		prev:    p,
		command: cmd,
	}
}
func (p *Preset) Concat(cmd ...Command) *Preset {
	preset := p
	for k := range cmd {
		preset = preset.concatCommand(cmd[k])
	}
	return preset
}

func (p *Preset) Exec(ctx *Context) (*Context, error) {
	if p == nil {
		return ctx, nil
	}
	c, err := p.prev.Exec(ctx)
	if err != nil {
		return nil, err
	}
	return p.command.Exec(c)
}
func NewPreset() *Preset {
	return nil
}
