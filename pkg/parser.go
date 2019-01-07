package uniris

import (
	"fmt"
)

type parser struct {
	tokens  []token
	current int
}

func (p *parser) parse() ([]statement, error) {
	statements := make([]statement, 0)
	for !p.isAtEnd() {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}
	return statements, nil
}

func (p *parser) statement() (statement, error) {
	if p.match(TokenFunction) {
		return p.functionStatement()
	}
	if p.match(TokenFor) {
		return p.forStatement()
	}
	if p.match(TokenIf) {
		return p.ifStatement()
	}
	if p.match(TokenPrint) {
		return p.printStatement()
	}
	if p.match(TokenReturn) {
		return p.returnStatement()
	}
	if p.match(TokenWhile) {
		return p.whileStatement()
	}
	if p.match(TokenLeftBracket) {
		return p.blockStatements()
	}
	return p.expressionStatement()
}

func (p *parser) returnStatement() (statement, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	return returnStatement{
		value: value,
	}, nil
}

func (p *parser) functionStatement() (statement, error) {
	name, err := p.consume(TokenIdentifier, "Expect function name")
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(TokenLeftParenthesis, "Expect '(' after function name"); err != nil {
		return nil, err
	}
	params := make([]token, 0)
	if !p.check(TokenRightParenthesis) {
		for {
			token, err := p.consume(TokenIdentifier, "Expect parameter name")
			if err != nil {
				return nil, err
			}
			params = append(params, token)
			if !p.match(TokenComma) {
				break
			}
		}
	}
	if _, err := p.consume(TokenRightParenthesis, "Expect ')' after parameters"); err != nil {
		return nil, err
	}
	if _, err := p.consume(TokenLeftBracket, "Expect '{' before function body"); err != nil {
		return nil, err
	}

	body, err := p.blockStatements()
	if err != nil {
		return nil, err
	}
	return funcStatement{
		body:   body.(blockStmt),
		name:   name,
		params: params,
	}, nil
}

func (p *parser) forStatement() (statement, error) {

	var init statement
	if p.check(TokenSemiColon) {
		init = nil
	} else if p.check(TokenIdentifier) {
		exp, err := p.assignement()
		if err != nil {
			return nil, err
		}
		init = exp
	} else {
		exp, err := p.expressionStatement()
		if err != nil {
			return nil, err
		}
		init = exp
	}

	p.advance()

	var cond expression
	if !p.check(TokenSemiColon) {
		exp, err := p.expression()
		if err != nil {
			return nil, err
		}
		cond = exp
	}

	if _, err := p.consume(TokenSemiColon, "Expected ; after loop condition"); err != nil {
		return nil, err
	}

	increment, err := p.expression()
	if err != nil {
		return nil, err
	}
	body, err := p.statement()
	if err != nil {
		return nil, err
	}
	if increment != nil {
		body = blockStmt{
			statements: []statement{
				body,
				expressionStmt{
					exp: increment,
				},
			},
		}
	}

	if cond == nil {
		cond = literalExpression{
			value: true,
		}
	}
	body = whileStatement{body: body, cond: cond}

	if init != nil {
		body = blockStmt{
			statements: []statement{
				init,
				body,
			},
		}
	}

	return body, nil

}

func (p *parser) whileStatement() (statement, error) {
	cond, err := p.expression()
	if err != nil {
		return nil, err
	}
	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return whileStatement{
		cond: cond,
		body: body,
	}, nil
}

func (p *parser) ifStatement() (statement, error) {
	cond, err := p.expression()
	if err != nil {
		return nil, err
	}

	thenStmt, err := p.statement()
	if err != nil {
		return nil, err
	}
	var elseStmt statement

	if p.match(TokenElse) {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		elseStmt = stmt
	}

	return ifStatement{
		cond:     cond,
		thenStmt: thenStmt,
		elseStmt: elseStmt,
	}, nil
}

func (p *parser) blockStatements() (statement, error) {
	statements := make([]statement, 0)
	for !p.check(TokenRightBracket) && !p.isAtEnd() {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}
	if _, err := p.consume(TokenRightBracket, "Expect } after block"); err != nil {
		return nil, err
	}
	return blockStmt{
		statements: statements,
	}, nil
}

func (p *parser) printStatement() (statement, error) {
	val, err := p.expression()
	if err != nil {
		return nil, err
	}
	return printStmt{exp: val}, nil
}

func (p *parser) expressionStatement() (statement, error) {
	exp, err := p.expression()
	if err != nil {
		return nil, err
	}
	return expressionStmt{exp}, nil
}

func (p *parser) expression() (expression, error) {
	return p.assignement()
}

func (p *parser) assignement() (expression, error) {
	exp, err := p.or()
	if err != nil {
		return nil, err
	}
	if p.match(TokenEqual) {
		eq := p.previous()
		val, err := p.assignement()
		if err != nil {
			return nil, err
		}

		return assignExpression{
			op:  eq,
			exp: val,
		}, nil
	}

	return exp, nil
}

func (p *parser) or() (expression, error) {
	exp, err := p.and()
	if err != nil {
		return nil, err
	}
	for p.match(TokenOr) {
		op := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		exp = logicalExpression{
			left:  exp,
			right: right,
			op:    op,
		}
	}
	return exp, nil
}

func (p *parser) and() (expression, error) {
	exp, err := p.equality()
	if err != nil {
		return nil, err
	}
	for p.match(TokenAnd) {
		op := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		exp = logicalExpression{
			op:    op,
			left:  exp,
			right: right,
		}
	}
	return exp, nil
}

func (p *parser) equality() (expression, error) {
	exp, err := p.comparison()
	if err != nil {
		return nil, err
	}
	for p.match(TokenBangEqual, TokenEqualEqual) {
		op := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		exp = binaryExpression{
			left:  exp,
			op:    op,
			right: right,
		}
	}

	return exp, nil
}

func (p *parser) comparison() (expression, error) {
	exp, err := p.addition()
	if err != nil {
		return nil, err
	}
	for p.match(TokenGreater, TokenGreaterEqual, TokenLess, TokenLessEqual) {
		op := p.previous()
		right, err := p.addition()
		if err != nil {
			return nil, err
		}
		exp = binaryExpression{
			left:  exp,
			op:    op,
			right: right,
		}
	}
	return exp, nil
}

func (p *parser) addition() (expression, error) {
	exp, err := p.multiplication()
	if err != nil {
		return nil, err
	}
	for p.match(TokenMinus, TokenPlus) {
		op := p.previous()
		right, err := p.multiplication()
		if err != nil {
			return nil, err
		}
		exp = binaryExpression{
			left:  exp,
			op:    op,
			right: right,
		}
	}
	return exp, nil
}

func (p *parser) multiplication() (expression, error) {
	exp, err := p.unary()
	if err != nil {
		return nil, err
	}
	for p.match(TokenSlash, TokenStar) {
		op := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		exp = binaryExpression{
			left:  exp,
			op:    op,
			right: right,
		}
	}
	return exp, nil
}

func (p *parser) unary() (expression, error) {
	if p.match(TokenBang, TokenMinus) {
		op := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return unaryExpression{
			op:    op,
			right: right,
		}, nil
	}

	return p.call()
}

func (p *parser) call() (expression, error) {
	exp, err := p.primary()
	if err != nil {
		return nil, err
	}
	for true {
		if p.match(TokenLeftParenthesis) {
			exp, err = p.finishCall(exp)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return exp, nil
}

func (p *parser) finishCall(callee expression) (expression, error) {
	args := make([]expression, 0)
	if !p.check(TokenRightParenthesis) {
		for {
			exp, err := p.expression()
			if err != nil {
				return nil, err
			}
			args = append(args, exp)
			if !p.match(TokenComma) {
				break
			}
		}
	}

	paren, err := p.consume(TokenRightParenthesis, "Expected ')' after arguments")
	if err != nil {
		return nil, err
	}
	return callExpression{
		args:   args,
		callee: callee,
		paren:  paren,
	}, nil
}

func (p *parser) primary() (expression, error) {
	if p.match(TokenFalse) {
		return literalExpression{value: false}, nil
	}
	if p.match(TokenTrue) {
		return literalExpression{value: true}, nil
	}
	if p.match(TokenNumber, TokenString) {
		return literalExpression{value: p.previous().Literal}, nil
	}
	if p.match(TokenIdentifier) {
		op := p.previous()
		if p.match(TokenEqual) {
			exp, err := p.expression()
			if err != nil {
				return nil, err
			}
			return assignExpression{
				op:  op,
				exp: exp,
			}, nil
		}
		return variableExpression{
			op: op,
		}, nil
	}
	if p.match(TokenLeftParenthesis) {
		exp, err := p.expression()
		if err != nil {
			return nil, err
		}
		if _, err := p.consume(TokenRightParenthesis, "Expect ')' after expression"); err != nil {
			return nil, err
		}
		return groupingExpression{exp: exp}, nil
	}

	err := p.error(p.peek(), "Expected expression")

	return nil, err
}

func (p *parser) match(ts ...TokenType) bool {
	for _, t := range ts {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *parser) consume(t TokenType, message string) (token, error) {
	if p.check(t) {
		return p.advance(), nil
	}

	err := p.error(p.peek(), message)
	return token{}, err
}

func (p *parser) error(tok token, message string) error {
	if tok.Type == TokenEndOfFile {
		return fmt.Errorf("Parsing error at end of line %d - %s", tok.Line, message)
	}
	return fmt.Errorf("Parsing error at %s of line %d - %s", tok.Lexeme, tok.Line, message)
}

func (p *parser) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *parser) advance() token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *parser) isAtEnd() bool {
	return p.peek().Type == TokenEndOfFile
}

func (p *parser) peek() token {
	return p.tokens[p.current]
}

func (p *parser) previous() token {
	return p.tokens[p.current-1]
}
