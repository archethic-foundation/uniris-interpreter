package main

import (
	"fmt"
)

type environment struct {
	enclosing *environment
	values    map[string]interface{}
}

func newGlobalEnvironment() *environment {
	globals := &environment{}
	globals.set("now", currentTimestampFunc{})

	return &environment{
		enclosing: globals,
	}
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

func (env *environment) get(name string) interface{} {
	v, exist := env.values[name]
	if exist {
		return v
	}

	if env.enclosing != nil {
		return env.enclosing.get(name)
	}

	panic(fmt.Sprintf("Undefined variable %s", name))
}
