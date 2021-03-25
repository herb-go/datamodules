package ncache

import (
	"bytes"
	"io"
)

type nextPath struct {
	next *nextPath
	*Path
}

func (n *nextPath) writeSelf(w io.Writer) (int, error) {
	if n.Path == nil {
		return 0, nil
	}
	pl, err := writeTokenAndData(w, tokenBeforeGroup, n.group)
	if err != nil {
		return 0, err
	}
	nl, err := writeTokenAndData(w, tokenBeforePath, n.name)
	if err != nil {
		return 0, err
	}
	return pl + nl, nil
}

func (n *nextPath) writeAll(w io.Writer) (int, error) {
	var written int
	np := n
	for np != nil {
		l, err := np.writeSelf(w)
		if err != nil {
			return 0, err
		}
		written += l
		np = np.next
	}
	return written, nil
}

type Path struct {
	prev  *Path
	name  []byte
	group []byte
}

func (p *Path) toNextPath(next *nextPath) *nextPath {
	np := &nextPath{
		next: next,
		Path: p,
	}
	if p == nil {
		return np
	}
	return p.prev.toNextPath(np)
}
func (p *Path) WriteTo(w io.Writer) (int64, error) {
	l, err := p.toNextPath(nil).writeAll(w)
	if err != nil {
		return 0, err
	}
	return int64(l), nil
}
func (p *Path) Append(group []byte, name []byte) *Path {
	return &Path{
		prev:  p,
		name:  name,
		group: group,
	}
}
func (p *Path) Equal(dst *Path) bool {
	if p == nil || dst == nil {
		return p == dst
	}

	if !bytes.Equal(p.name, dst.name) {
		return false
	}
	return p.prev.Equal(dst.prev)

}

func NewPath() *Path {
	return nil
}
