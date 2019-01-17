package uniris

import (
	"errors"
	"fmt"
	"reflect"
)

type expression interface {
	evaluate(*Environment) (interface{}, error)
}

//Variable assignation
type assignExpression struct {
	op  token
	exp expression
}

func (e assignExpression) evaluate(env *Environment) (interface{}, error) {
	value, err := e.exp.evaluate(env)
	if err != nil {
		return nil, err
	}

	env.Set(e.op.Lexeme, value)
	return nil, nil
}

//Variable execution
type variableExpression struct {
	op token
}

func (e variableExpression) evaluate(env *Environment) (interface{}, error) {
	return env.Get(e.op.Lexeme)
}

//Arithmetic (+ - * /) and logic (== !=  > < >= <=)
type binaryExpression struct {
	left  expression
	right expression
	op    token
}

func (e binaryExpression) evaluate(env *Environment) (interface{}, error) {
	left, err := e.left.evaluate(env)
	if err != nil {
		return nil, err
	}
	right, err := e.right.evaluate(env)
	if err != nil {
		return nil, err
	}

	switch e.op.Type {
	case TokenMinus:
		return left.(float64) - right.(float64), nil
	case TokenSlash:
		return left.(float64) / right.(float64), nil
	case TokenStar:
		return left.(float64) * right.(float64), nil
	case TokenPlus:
		if reflect.TypeOf(left).String() == "float64" && reflect.TypeOf(right).String() == "float64" {
			return left.(float64) + right.(float64), nil
		}
		return fmt.Sprintf("%v%v", left, right), nil
	case TokenGreater:
		return left.(float64) > right.(float64), nil
	case TokenGreaterEqual:
		return left.(float64) >= right.(float64), nil
	case TokenLess:
		return left.(float64) < right.(float64), nil
	case TokenLessEqual:
		return left.(float64) <= right.(float64), nil
	case TokenEqualEqual:
		return left == right, nil
	case TokenBangEqual:
		return left != right, nil
	default:
		return nil, errors.New("Not supported as binary expression")
	}
}

//Parenthesis and brackets
type groupingExpression struct {
	exp expression
}

func (e groupingExpression) evaluate(env *Environment) (interface{}, error) {
	return e.exp.evaluate(env)
}

//Not expression or negative one
type unaryExpression struct {
	op    token
	right expression
}

func (e unaryExpression) evaluate(env *Environment) (interface{}, error) {
	right, err := e.right.evaluate(env)
	if err != nil {
		return nil, err
	}
	switch e.op.Type {
	case TokenBang:
		return !isTruthy(right), nil
	case TokenMinus:
		return -right.(float64), nil
	}

	return nil, nil
}

//Number, string, booleans
type literalExpression struct {
	value interface{}
}

func (e literalExpression) evaluate(env *Environment) (interface{}, error) {
	return e.value, nil
}

//And, OR
type logicalExpression struct {
	left  expression
	op    token
	right expression
}

func (e logicalExpression) evaluate(env *Environment) (interface{}, error) {
	left, err := e.left.evaluate(env)
	if err != nil {
		return nil, err
	}
	if e.op.Type == TokenOr {
		if isTruthy(left) {
			return left, nil
		}
	} else {
		if !isTruthy(left) {
			return left, nil
		}
	}
	return e.right.evaluate(env)
}

type callExpression struct {
	callee expression
	paren  token
	args   []expression
}

func (e callExpression) evaluate(env *Environment) (interface{}, error) {
	callee, err := e.callee.evaluate(env)
	if err != nil {
		return nil, err
	}
	switch callee.(type) {
	case callable:
		args := make([]interface{}, 0)
		for _, arg := range e.args {
			val, err := arg.evaluate(env)
			if err != nil {
				return nil, err
			}
			args = append(args, val)
		}
		f := callee.(callable)
		return f.call(env, args...)
	default:
		return nil, errors.New("Can only call functions")
	}
}
