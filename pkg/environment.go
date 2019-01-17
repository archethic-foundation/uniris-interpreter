package uniris

import (
	"fmt"
)

//Environment contains the values and inner values storage for the interpreter context (variables, functions)
type Environment struct {
	enclosing *Environment
	values    map[string]interface{}
}

//NewEnvironment creates a new interpreter environment
func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		values:    make(map[string]interface{}, 0),
		enclosing: enclosing,
	}
}

func (env *Environment) Set(name string, value interface{}) {
	if env.enclosing != nil {
		_, err := env.enclosing.Get(name)
		if err != nil {
			if err.Error() == fmt.Sprintf("Undefined variable %s", name) {
				env.values[name] = value
				return
			}
			panic(err)
		}
		env.enclosing.Set(name, value)
	} else {
		env.values[name] = value
	}
}

func (env *Environment) Get(name string) (interface{}, error) {
	v, exist := env.values[name]
	if exist {
		return v, nil
	}

	if env.enclosing != nil {
		return env.enclosing.Get(name)
	}

	return nil, fmt.Errorf("Undefined variable %s", name)
}
