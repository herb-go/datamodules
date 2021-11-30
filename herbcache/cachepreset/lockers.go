package cachepreset

import (
	"github.com/herb-go/misc/waitingmap"
)

type Lockers struct {
	m *waitingmap.LockerMap
}

func (l *Lockers) Exec(preset *Preset) (newpreset *Preset, err error) {
	return preset.OverrideLockers(l), nil
}

func NewLockers() *Lockers {
	return &Lockers{
		m: waitingmap.NewLockerMap(),
	}
}
