package builderconfig

import (
	"fmt"

	"github.com/herb-go/datamodules/ncache"
)

type BuildConfig struct {
	Directives []*Directive
}

func (c *BuildConfig) CreateBuilders() ([]ncache.Builder, error) {
	result := []ncache.Builder{}
	for _, v := range c.Directives {
		b, err := v.CreateBuilder()
		if err != nil {
			return nil, err
		}
		result = append(result, b)
	}
	return result, nil
}

type Directive struct {
	Type   string
	Config func(v interface{}) error `config:", lazyload"`
}

func (d *Directive) CreateBuilder() (ncache.Builder, error) {
	switch d.Type {

	}
	return nil, fmt.Errorf("ncenter builderconfig:unknown directive type [%s]", d.Type)
}
