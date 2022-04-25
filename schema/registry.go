package schema

import (
	"fmt"
	"sync"

	"incr.app/incr/common"
)

var DefaultRegistry Registry

type Registry struct {
	l        *sync.RWMutex
	registry map[string]common.Schema
}

func (r *Registry) Lookup(name string) (common.Schema, bool) {
	r.l.RLock()
	defer r.l.RUnlock()
	schema, ok := r.registry[name]
	return schema, ok
}

func (r *Registry) Schemas() []string {
	result := make([]string, 0, len(r.registry))
	for name := range r.registry {
		result = append(result, name)
	}
	return result
}

func (r *Registry) Register(name string, schema common.Schema) error {
	if _, exist := r.registry[name]; exist {
		return fmt.Errorf("cannot replace existing schema %s", name)
	}
	r.registry[name] = schema
	return nil
}

func NewRegistry() Registry {
	r := Registry{
		l:        &sync.RWMutex{},
		registry: make(map[string]common.Schema, 2),
	}
	return r
}

func init() {
	DefaultRegistry = NewRegistry()
}
