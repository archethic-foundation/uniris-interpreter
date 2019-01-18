package uniris

import (
	"errors"
	"time"
)

type callable interface {
	call(*Environment, ...interface{}) (interface{}, error)
}

type function struct {
	declaration funcStatement
}

func (f function) call(env *Environment, args ...interface{}) (res interface{}, err error) {
	newEnvironment := NewEnvironment(env)

	if len(args) != len(f.declaration.params) {
		return nil, errors.New("Missing function parameters")
	}

	for i := 0; i < len(f.declaration.params); i++ {
		newEnvironment.Set(f.declaration.params[i].Lexeme, args[i])
	}

	defer func() {
		if x := recover(); x != nil {
			res = x
		}
	}()
	res, err = f.declaration.body.evaluate(newEnvironment)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//GLOBAL FUNCTIONS (BUILT-IN)
type currentTimestampFunc struct{}

func (f currentTimestampFunc) call(env *Environment, args ...interface{}) (interface{}, error) {
	return time.Now().Unix(), nil
}
