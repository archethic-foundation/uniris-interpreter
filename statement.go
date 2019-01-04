package main

import (
	"log"
)

type statement interface {
	evaluate(env *environment) interface{}
	print() string
}

type expressionStmt struct {
	exp expression
}

func (stmt expressionStmt) evaluate(env *environment) interface{} {
	return stmt.exp.evaluate(env)
}

func (stmt expressionStmt) print() string {
	return "expression"
}

type printStmt struct {
	exp expression
}

func (stmt printStmt) evaluate(env *environment) interface{} {
	value := stmt.exp.evaluate(env)
	log.Printf("%v", value)
	return nil
}

func (stmt printStmt) print() string {
	return "print"
}

type blockStmt struct {
	statements []statement
}

func (stmt blockStmt) evaluate(env *environment) interface{} {
	newenvironment := &environment{enclosing: env}

	for _, st := range stmt.statements {
		switch st.(type) {
		case returnStatement:
			if val := st.evaluate(newenvironment); val != nil {
				return val
			}
		default:
			st.evaluate(newenvironment)
		}
	}

	return nil
}

func (stmt blockStmt) print() string {
	return "block"
}

type ifStatement struct {
	cond     expression
	thenStmt statement
	elseStmt statement
}

func (stmt ifStatement) evaluate(env *environment) interface{} {
	cond := stmt.cond.evaluate(env)
	if isTruthy(cond) {
		stmt.thenStmt.evaluate(env)
	} else {
		if stmt.elseStmt != nil {
			elseStmt := stmt.elseStmt
			elseStmt.evaluate(env)
		}
	}
	return nil
}

func (stmt ifStatement) print() string {
	return "if"
}

type whileStatement struct {
	cond expression
	body statement
}

func (stmt whileStatement) evaluate(env *environment) interface{} {
	for isTruthy(stmt.cond.evaluate(env)) {
		stmt.body.evaluate(env)
	}
	return nil
}

func (stmt whileStatement) print() string {
	return "while"
}

type funcStatement struct {
	name   token
	params []token
	body   blockStmt
}

func (stmt funcStatement) print() string {
	return "func"
}

func (stmt funcStatement) evaluate(env *environment) interface{} {
	f := function{
		declaration: stmt,
	}
	env.set(stmt.name.Lexeme, f)
	return nil
}

type returnStatement struct {
	value expression
}

func (stmt returnStatement) evaluate(env *environment) interface{} {
	value := stmt.value.evaluate(env)
	if value != nil {
		//Back to the top of the stack on the call statement
		//To ensure this, we panic to make like handled exception
		panic(value)
	}
	return nil
}

func (stmt returnStatement) print() string {
	return "return"
}

func isTruthy(val interface{}) bool {
	if val == nil {
		return false
	}
	switch val.(type) {
	case bool:
		return val.(bool)
	}
	return true
}
