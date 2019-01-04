package main

import (
	"time"
)

type callable interface {
	call(*environment, ...interface{}) interface{}
}

type function struct {
	declaration funcStatement
}

func (f function) call(env *environment, args ...interface{}) (res interface{}) {
	newenvironment := &environment{enclosing: env}
	for i := 0; i < len(f.declaration.params); i++ {
		newenvironment.set(f.declaration.params[i].Lexeme, args[i])
	}

	defer func() {
		if x := recover(); x != nil {
			res = x
		}
	}()
	res = f.declaration.body.evaluate(newenvironment)
	return res
}

//GLOBAL FUNCTIONS (BUILT-IN)
type currentTimestampFunc struct{}

func (f currentTimestampFunc) call(env *environment, args ...interface{}) interface{} {
	return time.Now().Unix()
}
