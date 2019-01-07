package uniris

import (
	"fmt"
)

type environment struct {
	enclosing *environment
	values    map[string]interface{}
}

func (env *environment) set(name string, value interface{}) {
	if env.values == nil {
		env.values = make(map[string]interface{})
	}
	if env.enclosing != nil {
		if _, exist := env.enclosing.values[name]; exist {
			env.enclosing.set(name, value)
		} else {
			env.values[name] = value
		}
	} else {
		env.values[name] = value
	}
}

func (env *environment) get(name string) (interface{}, error) {
	v, exist := env.values[name]
	if exist {
		return v, nil
	}

	if env.enclosing != nil {
		return env.enclosing.get(name)
	}

	return nil, fmt.Errorf("Undefined variable %s", name)
}
