package uniris

import (
	"fmt"
)

type statement interface {
	evaluate(env *Environment) (interface{}, error)
}

type expressionStmt struct {
	exp expression
}

func (stmt expressionStmt) evaluate(env *Environment) (interface{}, error) {
	return stmt.exp.evaluate(env)
}

type printStmt struct {
	exp expression
}

func (stmt printStmt) evaluate(env *Environment) (interface{}, error) {
	value, err := stmt.exp.evaluate(env)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v\n", value)
	return nil, nil
}

type blockStmt struct {
	statements []statement
}

func (stmt blockStmt) evaluate(env *Environment) (interface{}, error) {
	newenvironment := NewEnvironment(env)

	for _, st := range stmt.statements {
		switch st.(type) {
		case returnStatement:
			val, err := st.evaluate(newenvironment)
			if err != nil {
				return nil, err
			}
			return val, nil
		default:
			if _, err := st.evaluate(newenvironment); err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

type ifStatement struct {
	cond     expression
	thenStmt statement
	elseStmt statement
}

func (stmt ifStatement) evaluate(env *Environment) (interface{}, error) {
	cond, err := stmt.cond.evaluate(env)
	if err != nil {
		return nil, err
	}
	if isTruthy(cond) {
		if _, err := stmt.thenStmt.evaluate(env); err != nil {
			return nil, err
		}
	} else {
		if stmt.elseStmt != nil {
			elseStmt := stmt.elseStmt
			elseStmt.evaluate(env)
		}
	}
	return nil, nil
}

type whileStatement struct {
	cond expression
	body statement
}

func (stmt whileStatement) evaluate(env *Environment) (interface{}, error) {
	for {
		val, err := stmt.cond.evaluate(env)
		if err != nil {
			return nil, err
		}
		if !isTruthy(val) {
			break
		}
		stmt.body.evaluate(env)
	}
	return nil, nil
}

type funcStatement struct {
	name   token
	params []token
	body   blockStmt
}

func (stmt funcStatement) evaluate(env *Environment) (interface{}, error) {
	f := function{
		declaration: stmt,
	}
	env.Set(stmt.name.Lexeme, f)
	return nil, nil
}

type returnStatement struct {
	value expression
}

func (stmt returnStatement) evaluate(env *Environment) (interface{}, error) {
	value, err := stmt.value.evaluate(env)
	if err != nil {
		return nil, err
	}
	if value != nil {
		//Back to the top of the stack on the call statement
		//To ensure this, we panic to make like handled exception
		panic(value)
	}
	return nil, nil
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
