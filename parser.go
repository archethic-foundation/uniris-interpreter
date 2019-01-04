package main

import (
	"fmt"
)

type parser struct {
	tokens  []token
	current int
}

func (p *parser) parse() []statement {
	statements := make([]statement, 0)
	for !p.isAtEnd() {
		statements = append(statements, p.statement())
	}
	return statements
}

func (p *parser) statement() statement {
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
		return blockStmt{statements: p.blockStatements()}
	}
	return p.expressionStatement()
}

func (p *parser) returnStatement() statement {
	value := p.expression()
	return returnStatement{
		value: value,
	}
}

func (p *parser) functionStatement() statement {
	name := p.consume(TokenIdentifier, "Expect function name")
	p.consume(TokenLeftParenthesis, "Expect '(' after function name")
	params := make([]token, 0)
	if !p.check(TokenRightParenthesis) {
		for {
			params = append(params, p.consume(TokenIdentifier, "Expect parameter name"))
			if !p.match(TokenComma) {
				break
			}
		}
	}
	p.consume(TokenRightParenthesis, "Expect ')' after parameters")
	p.consume(TokenLeftBracket, "Expect '{' before function body")
	body := blockStmt{
		statements: p.blockStatements(),
	}
	return funcStatement{
		body:   body,
		name:   name,
		params: params,
	}
}

func (p *parser) forStatement() statement {

	var init statement
	if p.check(TokenSemiColon) {
		init = nil
	} else if p.check(TokenIdentifier) {
		init = p.assignement()
	} else {
		init = p.expressionStatement()
	}

	p.advance()

	var cond expression
	if !p.check(TokenSemiColon) {
		cond = p.expression()
	}

	p.consume(TokenSemiColon, "Expected ; after loop condition")

	increment := p.expression()
	body := p.statement()
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

	return body

}

func (p *parser) whileStatement() statement {
	cond := p.expression()
	body := p.statement()

	return whileStatement{
		cond: cond,
		body: body,
	}
}

func (p *parser) ifStatement() statement {
	cond := p.expression()

	thenStmt := p.statement()
	var elseStmt statement

	if p.match(TokenElse) {
		elseStmt = p.statement()
	}

	return ifStatement{
		cond:     cond,
		thenStmt: thenStmt,
		elseStmt: elseStmt,
	}
}

func (p *parser) blockStatements() []statement {
	statements := make([]statement, 0)
	for !p.check(TokenRightBracket) && !p.isAtEnd() {
		statements = append(statements, p.statement())
	}
	p.consume(TokenRightBracket, "Expect } after block")
	return statements
}

func (p *parser) printStatement() statement {
	val := p.expression()
	return printStmt{exp: val}
}

func (p *parser) expressionStatement() statement {
	exp := p.expression()
	return expressionStmt{exp}
}

func (p *parser) expression() expression {
	return p.assignement()
}

func (p *parser) assignement() expression {
	exp := p.or()
	if p.match(TokenEqual) {
		eq := p.previous()
		val := p.assignement()

		return assignExpression{
			op:  eq,
			exp: val,
		}
	}

	return exp
}

func (p *parser) or() expression {
	exp := p.and()
	for p.match(TokenOr) {
		op := p.previous()
		right := p.and()
		exp = logicalExpression{
			left:  exp,
			right: right,
			op:    op,
		}
	}
	return exp
}

func (p *parser) and() expression {
	exp := p.equality()
	for p.match(TokenAnd) {
		op := p.previous()
		right := p.equality()
		exp = logicalExpression{
			op:    op,
			left:  exp,
			right: right,
		}
	}
	return exp
}

func (p *parser) equality() expression {
	exp := p.comparison()
	for p.match(TokenBangEqual, TokenEqualEqual) {
		op := p.previous()
		right := p.comparison()
		exp = binaryExpression{
			left:  exp,
			op:    op,
			right: right,
		}
	}

	return exp
}

func (p *parser) comparison() expression {
	exp := p.addition()
	for p.match(TokenGreater, TokenGreaterEqual, TokenLess, TokenLessEqual) {
		op := p.previous()
		right := p.addition()
		exp = binaryExpression{
			left:  exp,
			op:    op,
			right: right,
		}
	}
	return exp
}

func (p *parser) addition() expression {
	exp := p.multiplication()
	for p.match(TokenMinus, TokenPlus) {
		op := p.previous()
		right := p.multiplication()
		exp = binaryExpression{
			left:  exp,
			op:    op,
			right: right,
		}
	}
	return exp
}

func (p *parser) multiplication() expression {
	exp := p.unary()
	for p.match(TokenSlash, TokenStar) {
		op := p.previous()
		right := p.unary()
		exp = binaryExpression{
			left:  exp,
			op:    op,
			right: right,
		}
	}
	return exp
}

func (p *parser) unary() expression {
	if p.match(TokenBang, TokenMinus) {
		op := p.previous()
		right := p.unary()
		return unaryExpression{
			op:    op,
			right: right,
		}
	}

	return p.call()
}

func (p *parser) call() expression {
	exp := p.primary()
	for true {
		if p.match(TokenLeftParenthesis) {
			exp = p.finishCall(exp)
		} else {
			break
		}
	}

	return exp
}

func (p *parser) finishCall(callee expression) expression {
	args := make([]expression, 0)
	if !p.check(TokenRightParenthesis) {
		for {
			args = append(args, p.expression())
			if !p.match(TokenComma) {
				break
			}
		}
	}

	paren := p.consume(TokenRightParenthesis, "Expected ')' after arguments")
	return callExpression{
		args:   args,
		callee: callee,
		paren:  paren,
	}
}

func (p *parser) primary() expression {
	if p.match(TokenFalse) {
		return literalExpression{value: false}
	}
	if p.match(TokenTrue) {
		return literalExpression{value: true}
	}
	if p.match(TokenNumber, TokenString) {
		return literalExpression{value: p.previous().Literal}
	}
	if p.match(TokenIdentifier) {
		op := p.previous()
		if p.match(TokenEqual) {
			exp := p.expression()
			return assignExpression{
				op:  op,
				exp: exp,
			}
		}
		return variableExpression{
			op: op,
		}
	}
	if p.match(TokenLeftParenthesis) {
		exp := p.expression()
		p.consume(TokenRightParenthesis, "Expect ')' after expression")
		return groupingExpression{
			exp: exp,
		}
	}

	p.error(p.peek(), "Expected expression")

	return nil
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

func (p *parser) consume(t TokenType, message string) token {
	if p.check(t) {
		return p.advance()
	}

	p.error(p.peek(), message)
	return token{}
}

func (p *parser) error(tok token, message string) {
	if tok.Type == TokenEndOfFile {
		panic(fmt.Sprintf("Parsing error at end of line %d - %s", tok.Line, message))
	} else {
		panic(fmt.Sprintf("Parsing error at %s of line %d - %s", tok.Lexeme, tok.Line, message))
	}
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
