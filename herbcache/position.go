package herbcache

import (
	"bytes"
)

type Position struct {
	Parent *Position
	Name   []byte
	Group  []byte
}

func (p *Position) Equal(dst *Position) bool {
	if p == nil || dst == nil {
		return p == dst
	}

	if !bytes.Equal(p.Name, dst.Name) {
		return false
	}
	return p.Parent.Equal(dst.Parent)

}
func (p *Position) RootDirectory() *Directory {
	return p.toPath(nil)
}
func (p *Position) toPath(next *Directory) *Directory {
	d := &Directory{
		Next:     next,
		Position: p,
	}
	if p == nil {
		return d
	}
	return p.Parent.toPath(d)
}

func (p *Position) Append(group []byte, name []byte) *Position {
	return &Position{
		Parent: p,
		Name:   name,
		Group:  group,
	}
}

type Directory struct {
	*Position
	Next *Directory
}
