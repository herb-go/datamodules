package cachepreset

type Commands struct {
	prev    *Commands
	command Command
}

func (p *Commands) concatCommand(cmd Command) *Commands {
	return &Commands{
		prev:    p,
		command: cmd,
	}
}
func (p *Commands) Concat(cmd ...Command) *Commands {
	preset := p
	for k := range cmd {
		preset = preset.concatCommand(cmd[k])
	}
	return preset
}
func (p *Commands) Length() int {
	if p == nil {
		return 0
	}
	return p.prev.Length() + 1
}
func (p *Commands) Exec(preset *Preset) (*Preset, error) {
	if p == nil {
		return preset, nil
	}
	c, err := p.prev.Exec(preset)
	if err != nil {
		return nil, err
	}
	return p.command.Exec(c)
}
func NewCommands() *Commands {
	return nil
}
