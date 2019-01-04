package main

import (
	"fmt"
	"reflect"
)

type expression interface {
	evaluate(*environment) interface{}
	print() string
}

//Variable assignation
type assignExpression struct {
	op  token
	exp expression
}

func (e assignExpression) evaluate(env *environment) interface{} {
	value := e.exp.evaluate(env)
	env.set(e.op.Lexeme, value)
	return value
}

func (e assignExpression) print() string {
	return "assign variable"
}

//Variable execution
type variableExpression struct {
	op token
}

func (e variableExpression) evaluate(env *environment) interface{} {
	return env.get(e.op.Lexeme)
}

func (e variableExpression) print() string {
	return "execute variable"
}

//Arithmetic (+ - * /) and logic (== !=  > < >= <=)
type binaryExpression struct {
	left  expression
	right expression
	op    token
}

func (e binaryExpression) evaluate(env *environment) interface{} {
	left := e.left.evaluate(env)
	right := e.right.evaluate(env)

	switch e.op.Type {
	case TokenMinus:
		return left.(float64) - right.(float64)
	case TokenSlash:
		return left.(float64) / right.(float64)
	case TokenStar:
		return left.(float64) * right.(float64)
	case TokenPlus:
		if reflect.TypeOf(left).String() == "float64" && reflect.TypeOf(right).String() == "float64" {
			return left.(float64) + right.(float64)
		}
		return fmt.Sprintf("%v%v", left, right)
	case TokenGreater:
		return left.(float64) > right.(float64)
	case TokenGreaterEqual:
		return left.(float64) >= right.(float64)
	case TokenLess:
		return left.(float64) < right.(float64)
	case TokenLessEqual:
		return left.(float64) <= right.(float64)
	case TokenEqualEqual:
		return left == right
	case TokenBangEqual:
		return left != right
	}

	return nil
}

func (e binaryExpression) print() string {
	return "binary expression"
}

//Parenthesis and brackets
type groupingExpression struct {
	exp expression
}

func (e groupingExpression) evaluate(env *environment) interface{} {
	return e.exp.evaluate(env)
}

func (e groupingExpression) print() string {
	return "group expression"
}

//Not expression or negative one
type unaryExpression struct {
	op    token
	right expression
}

func (e unaryExpression) evaluate(env *environment) interface{} {
	right := e.right.evaluate(env)
	switch e.op.Type {
	case TokenBang:
		return !isTruthy(right)
	case TokenMinus:
		return -right.(float64)
	}

	return nil
}

func (e unaryExpression) print() string {
	return "unary expression"
}

//Number, string, booleans
type literalExpression struct {
	value interface{}
}

func (e literalExpression) evaluate(env *environment) interface{} {
	return e.value
}

func (e literalExpression) print() string {
	return "literal expression"
}

//And, OR
type logicalExpression struct {
	left  expression
	op    token
	right expression
}

func (e logicalExpression) evaluate(env *environment) interface{} {
	left := e.left.evaluate(env)
	if e.op.Type == TokenOr {
		if isTruthy(left) {
			return left
		}
	} else {
		if !isTruthy(left) {
			return left
		}
	}
	return e.right.evaluate(env)
}

func (e logicalExpression) print() string {
	return "logical expression"
}

type callExpression struct {
	callee expression
	paren  token
	args   []expression
}

func (e callExpression) evaluate(env *environment) interface{} {
	callee := e.callee.evaluate(env)
	switch callee.(type) {
	case callable:
		args := make([]interface{}, 0)
		for _, arg := range e.args {
			args = append(args, arg.evaluate(env))
		}
		f := callee.(callable)
		return f.call(env, args...)
	default:
		panic("Can only call functions")
	}
}

func (e callExpression) print() string {
	return "call expression"
}
