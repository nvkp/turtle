package turtle

import (
	"fmt"
	"reflect"

	"github.com/nvkp/turtle/graph"
	"github.com/nvkp/turtle/scanner"
)

type Config struct {
	Base     string
	Prefixes map[string]string
}

func (c *Config) Marshal(v interface{}) ([]byte, error) {
	g := graph.NewWithOptions(
		graph.Options{
			Base:     c.Base,
			Prefixes: c.Prefixes,
		})
	if err := marshal(g, reflect.ValueOf(v)); err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	return g.Bytes()
}

func (c *Config) Unmarshal(data []byte, v interface{}) error {
	if v == nil {
		return ErrNilValue
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return ErrNoPointerValue
	}

	err := unmarshal(
		scanner.NewWithOptions(data,
			scanner.Options{
				Base:     c.Base,
				Prefixes: c.Prefixes,
			}),
		rv)
	if err != nil {
		return fmt.Errorf("unmarshal: %v", err)
	}

	return nil
}
