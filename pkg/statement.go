package uniris

import "fmt"

type statement interface {
	evaluate(env *environment) (interface{}, error)
	print() string
}

type expressionStmt struct {
	exp expression
}

func (stmt expressionStmt) evaluate(env *environment) (interface{}, error) {
	return stmt.exp.evaluate(env)
}

func (stmt expressionStmt) print() string {
	return "expression"
}

type printStmt struct {
	exp expression
}

func (stmt printStmt) evaluate(env *environment) (interface{}, error) {
	value, err := stmt.exp.evaluate(env)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v\n", value)
	return nil, nil
}

func (stmt printStmt) print() string {
	return "print"
}

type blockStmt struct {
	statements []statement
}

func (stmt blockStmt) evaluate(env *environment) (interface{}, error) {
	newenvironment := &environment{enclosing: env}

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

func (stmt blockStmt) print() string {
	return "block"
}

type ifStatement struct {
	cond     expression
	thenStmt statement
	elseStmt statement
}

func (stmt ifStatement) evaluate(env *environment) (interface{}, error) {
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

func (stmt ifStatement) print() string {
	return "if"
}

type whileStatement struct {
	cond expression
	body statement
}

func (stmt whileStatement) evaluate(env *environment) (interface{}, error) {
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

func (stmt funcStatement) evaluate(env *environment) (interface{}, error) {
	f := function{
		declaration: stmt,
	}
	env.set(stmt.name.Lexeme, f)
	return nil, nil
}

type returnStatement struct {
	value expression
}

func (stmt returnStatement) evaluate(env *environment) (interface{}, error) {
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
