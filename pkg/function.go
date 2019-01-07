package uniris

import (
	"time"
)

type callable interface {
	call(*environment, ...interface{}) (interface{}, error)
}

type function struct {
	declaration funcStatement
}

func (f function) call(env *environment, args ...interface{}) (res interface{}, err error) {
	newenvironment := &environment{enclosing: env}
	for i := 0; i < len(f.declaration.params); i++ {
		newenvironment.set(f.declaration.params[i].Lexeme, args[i])
	}

	defer func() {
		if x := recover(); x != nil {
			res = x
		}
	}()
	res, err = f.declaration.body.evaluate(newenvironment)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//GLOBAL FUNCTIONS (BUILT-IN)
type currentTimestampFunc struct{}

func (f currentTimestampFunc) call(env *environment, args ...interface{}) (interface{}, error) {
	return time.Now().Unix(), nil
}
