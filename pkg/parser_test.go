package uniris

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParserPreviousToken(t *testing.T) {
	p := parser{
		tokens: []token{
			token{
				Type: TokenNumber,
			},
			token{
				Type: TokenPlus,
			},
			token{
				Type: TokenNumber,
			},
		},
	}

	p.current = 2

	assert.Equal(t, TokenPlus, p.previous().Type)
}

func TestParserPeek(t *testing.T) {
	p := parser{
		tokens: []token{
			token{
				Type: TokenNumber,
			},
		},
	}

	assert.Equal(t, TokenNumber, p.peek().Type)
}

func TestParserIsAtEnd(t *testing.T) {
	p := parser{
		tokens: []token{
			token{
				Type: TokenEndOfFile,
			},
		},
	}

	assert.True(t, p.isAtEnd())
}

func TestParserAdvance(t *testing.T) {
	p := parser{
		tokens: []token{
			token{
				Type: TokenNumber,
			},
			token{
				Type: TokenPlus,
			},
			token{
				Type: TokenNumber,
			},
		},
	}

	assert.Equal(t, TokenNumber, p.advance().Type)
	assert.Equal(t, TokenPlus, p.advance().Type)
	assert.Equal(t, TokenNumber, p.advance().Type)
}

func TestParserCheck(t *testing.T) {
	p := parser{
		tokens: []token{
			token{
				Type: TokenNumber,
			},
		},
	}

	assert.True(t, p.check(TokenNumber))
}

func TestParserError(t *testing.T) {
	p := parser{}
	assert.Error(t, p.error(token{Type: TokenEndOfFile, Line: 1}, "Invalid"), "Parsing error at end of line 1 - Invalid")
	assert.Error(t, p.error(token{Type: TokenNumber, Line: 1, Lexeme: "2"}, "Invalid"), "Parsing error at 2 of line 1 - Invalid")
}

func TestParserConsume(t *testing.T) {
	p := parser{
		tokens: []token{
			token{
				Type: TokenNumber,
			},
			token{
				Type: TokenPlus,
			},
			token{
				Type: TokenNumber,
			},
		},
	}

	token, err := p.consume(TokenNumber, "Invalid number")
	assert.Equal(t, TokenNumber, token.Type)
	assert.Nil(t, err)

	token, err = p.consume(TokenNumber, "Invalid Operator")
	assert.NotNil(t, err)
}

func TestParserMatch(t *testing.T) {
	p := parser{
		tokens: []token{
			token{
				Type: TokenNumber,
			},
			token{
				Type: TokenPlus,
			},
			token{
				Type: TokenNumber,
			},
		},
	}

	assert.True(t, p.match(TokenNumber, TokenPlus))

	assert.False(t, p.match(TokenIf))
}

func TestParserPrimaryLiteralExpression(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: TokenFalse},
			token{Type: TokenTrue},
			token{Type: TokenNumber, Literal: "10"},
			token{Type: TokenString, Literal: "hello"},
		},
	}

	exp, err := p.primary()
	assert.Nil(t, err)
	assert.Equal(t, literalExpression{value: false}, exp)

	exp, err = p.primary()
	assert.Nil(t, err)
	assert.Equal(t, literalExpression{value: true}, exp)

	exp, err = p.primary()
	assert.Nil(t, err)
	assert.Equal(t, literalExpression{value: "10"}, exp)

	exp, err = p.primary()
	assert.Nil(t, err)
	assert.Equal(t, literalExpression{value: "hello"}, exp)
}

func TestParserPrimaryAssignExpression(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: TokenIdentifier},
			token{Type: TokenEqual},
			token{Type: TokenNumber, Literal: 10},
			token{Type: TokenEndOfFile},
		},
	}

	exp, err := p.primary()
	assert.Nil(t, err)
	assert.Equal(t, assignExpression{
		op: token{
			Type: TokenIdentifier,
		},
		exp: literalExpression{
			value: 10,
		},
	}, exp)
}

func TestParserPrimaryVariableExpression(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: TokenIdentifier},
			token{Type: TokenEndOfFile},
		},
	}

	exp, err := p.primary()
	assert.Nil(t, err)
	assert.Equal(t, variableExpression{
		op: token{
			Type: TokenIdentifier,
		},
	}, exp)
}

func TestParserPrimaryGroupingExpression(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: TokenLeftParenthesis},
			token{Type: TokenTrue},
			token{Type: TokenRightParenthesis},
			token{Type: TokenEndOfFile},
		},
	}

	exp, err := p.primary()
	assert.Nil(t, err)
	assert.Equal(t, groupingExpression{
		exp: literalExpression{
			value: true,
		},
	}, exp)
}

func TestParserFinishCall(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: TokenNumber, Literal: 10},
			token{Type: TokenRightParenthesis},
			token{Type: TokenEndOfFile},
		},
	}
	exp, err := p.finishCall(testFuncExpression{})
	assert.Nil(t, err)
	assert.Equal(t, callExpression{
		args: []expression{
			literalExpression{
				value: 10,
			},
		},
		paren: token{
			Type: TokenRightParenthesis,
		},
		callee: testFuncExpression{},
	}, exp)
}

func TestParserCall(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: TokenIdentifier},
			token{Type: TokenLeftParenthesis},
			token{Type: TokenRightParenthesis},
			token{Type: TokenEndOfFile},
		},
	}

	exp, err := p.call()
	assert.Nil(t, err)
	assert.Equal(t, callExpression{
		args: []expression{},
		callee: variableExpression{
			op: token{Type: TokenIdentifier},
		},
		paren: token{Type: TokenRightParenthesis},
	}, exp)
}

func TestParserUnary(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: TokenBang},
			token{Type: TokenTrue},
			token{Type: TokenMinus},
			token{Type: TokenNumber, Literal: 10},
			token{Type: TokenEndOfFile},
		},
	}

	exp, err := p.unary()
	assert.Nil(t, err)
	assert.Equal(t, unaryExpression{
		op: token{Type: TokenBang},
		right: literalExpression{
			value: true,
		},
	}, exp)

	exp, err = p.unary()
	assert.Nil(t, err)
	assert.Equal(t, unaryExpression{
		op: token{Type: TokenMinus},
		right: literalExpression{
			value: 10,
		},
	}, exp)
}

func TestParserMultiplication(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: TokenNumber, Literal: 10},
			token{Type: TokenStar},
			token{Type: TokenNumber, Literal: 10},
			token{Type: TokenEndOfFile},
		},
	}

	exp, err := p.multiplication()
	assert.Nil(t, err)
	assert.Equal(t, binaryExpression{
		left:  literalExpression{value: 10},
		op:    token{Type: TokenStar},
		right: literalExpression{value: 10},
	}, exp)
}

func TestParserAddition(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: TokenNumber, Literal: 10},
			token{Type: TokenPlus},
			token{Type: TokenNumber, Literal: 10},
			token{Type: TokenNumber, Literal: 10},
			token{Type: TokenMinus},
			token{Type: TokenNumber, Literal: 10},
			token{Type: TokenEndOfFile},
		},
	}

	exp, err := p.addition()
	assert.Nil(t, err)
	assert.Equal(t, binaryExpression{
		left:  literalExpression{value: 10},
		op:    token{Type: TokenPlus},
		right: literalExpression{value: 10},
	}, exp)

	exp, err = p.addition()
	assert.Nil(t, err)
	assert.Equal(t, binaryExpression{
		left:  literalExpression{value: 10},
		op:    token{Type: TokenMinus},
		right: literalExpression{value: 10},
	}, exp)
}

func TestParserComparison(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: TokenNumber, Literal: 10},
			token{Type: TokenGreater},
			token{Type: TokenNumber, Literal: 10},
			token{Type: TokenNumber, Literal: 10},
			token{Type: TokenGreaterEqual},
			token{Type: TokenNumber, Literal: 10},
			token{Type: TokenNumber, Literal: 10},
			token{Type: TokenLess},
			token{Type: TokenNumber, Literal: 10},
			token{Type: TokenNumber, Literal: 10},
			token{Type: TokenLessEqual},
			token{Type: TokenNumber, Literal: 10},
			token{Type: TokenEndOfFile},
		},
	}

	exp, err := p.comparison()
	assert.Nil(t, err)
	assert.Equal(t, binaryExpression{
		left:  literalExpression{value: 10},
		op:    token{Type: TokenGreater},
		right: literalExpression{value: 10},
	}, exp)

	exp, err = p.comparison()
	assert.Nil(t, err)
	assert.Equal(t, binaryExpression{
		left:  literalExpression{value: 10},
		op:    token{Type: TokenGreaterEqual},
		right: literalExpression{value: 10},
	}, exp)

	exp, err = p.comparison()
	assert.Nil(t, err)
	assert.Equal(t, binaryExpression{
		left:  literalExpression{value: 10},
		op:    token{Type: TokenLess},
		right: literalExpression{value: 10},
	}, exp)

	exp, err = p.comparison()
	assert.Nil(t, err)
	assert.Equal(t, binaryExpression{
		left:  literalExpression{value: 10},
		op:    token{Type: TokenLessEqual},
		right: literalExpression{value: 10},
	}, exp)
}

func TestParserEquality(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: TokenNumber, Literal: 10},
			token{Type: TokenEqualEqual},
			token{Type: TokenNumber, Literal: 10},
			token{Type: TokenIdentifier},
			token{Type: TokenBangEqual},
			token{Type: TokenNumber, Literal: 10},
			token{Type: TokenEndOfFile},
		},
	}

	exp, err := p.equality()
	assert.Nil(t, err)
	assert.Equal(t, binaryExpression{
		left:  literalExpression{value: 10},
		op:    token{Type: TokenEqualEqual},
		right: literalExpression{value: 10},
	}, exp)

	exp, err = p.equality()
	assert.Nil(t, err)
	assert.Equal(t, binaryExpression{
		left:  variableExpression{op: token{Type: TokenIdentifier}},
		op:    token{Type: TokenBangEqual},
		right: literalExpression{value: 10},
	}, exp)
}

func TestParserAnd(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: TokenTrue},
			token{Type: TokenAnd},
			token{Type: TokenTrue},
			token{Type: TokenEndOfFile},
		},
	}

	exp, err := p.and()
	assert.Nil(t, err)
	assert.Equal(t, logicalExpression{
		left:  literalExpression{value: true},
		op:    token{Type: TokenAnd},
		right: literalExpression{value: true},
	}, exp)
}

func TestParserOr(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: TokenTrue},
			token{Type: TokenOr},
			token{Type: TokenTrue},
			token{Type: TokenEndOfFile},
		},
	}

	exp, err := p.or()
	assert.Nil(t, err)
	assert.Equal(t, logicalExpression{
		left:  literalExpression{value: true},
		op:    token{Type: TokenOr},
		right: literalExpression{value: true},
	}, exp)
}

func TestParserAssignement(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: TokenIdentifier},
			token{Type: TokenEqual},
			token{Type: TokenNumber, Literal: 10},
			token{Type: TokenEndOfFile},
		},
	}

	exp, err := p.assignement()
	assert.Nil(t, err)
	assert.Equal(t, assignExpression{
		exp: literalExpression{value: 10},
		op:  token{Type: TokenIdentifier},
	}, exp)
}

func TestParserExpressionStatement(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: TokenNumber, Literal: 10},
			token{Type: TokenEndOfFile},
		},
	}

	stmt, err := p.expressionStatement()
	assert.Nil(t, err)
	assert.Equal(t, expressionStmt{
		exp: literalExpression{value: 10},
	}, stmt)

}

func TestParserPrintStatement(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: TokenNumber, Literal: 10},
			token{Type: TokenEndOfFile},
		},
	}

	stmt, err := p.printStatement()
	assert.Nil(t, err)
	assert.Equal(t, printStmt{
		exp: literalExpression{value: 10},
	}, stmt)
}

func TestParserBlockStatement(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: TokenPrint},
			token{Type: TokenNumber, Literal: 10},
			token{Type: TokenRightBracket},
			token{Type: TokenEndOfFile},
		},
	}

	stmt, err := p.blockStatements()
	assert.Nil(t, err)
	assert.Equal(t, blockStmt{
		statements: []statement{
			printStmt{
				exp: literalExpression{value: 10},
			},
		},
	}, stmt)
}

func TestParserIfStatement(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: TokenTrue},
			token{Type: TokenNumber},
			token{Type: TokenElse},
			token{Type: TokenNumber},
			token{Type: TokenEndOfFile},
		},
	}

	stmt, err := p.ifStatement()
	assert.Nil(t, err)
	assert.Equal(t, ifStatement{
		cond:     literalExpression{value: true},
		thenStmt: expressionStmt{exp: literalExpression{}},
		elseStmt: expressionStmt{exp: literalExpression{}},
	}, stmt)
}

func TestParserWhileStatement(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: TokenTrue},
			token{Type: TokenNumber},
			token{Type: TokenEndOfFile},
		},
	}

	stmt, err := p.whileStatement()
	assert.Nil(t, err)
	assert.Equal(t, whileStatement{
		cond: literalExpression{value: true},
		body: expressionStmt{exp: literalExpression{}},
	}, stmt)
}

func TestParserForStatement(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: TokenIdentifier, Lexeme: "i"},
			token{Type: TokenEqual},
			token{Type: TokenNumber, Literal: 0},
			token{Type: TokenSemiColon},
			token{Type: TokenIdentifier, Lexeme: "i"},
			token{Type: TokenLess},
			token{Type: TokenNumber, Literal: 10},
			token{Type: TokenSemiColon},
			token{Type: TokenIdentifier, Lexeme: "i"},
			token{Type: TokenEqual},
			token{Type: TokenIdentifier, Lexeme: "i"},
			token{Type: TokenPlus},
			token{Type: TokenNumber, Literal: 1},
			token{Type: TokenPrint},
			token{Type: TokenIdentifier, Lexeme: "i"},
			token{Type: TokenEndOfFile},
		},
	}

	stmt, err := p.forStatement()
	assert.Nil(t, err)
	assert.Equal(t, blockStmt{
		statements: []statement{
			assignExpression{
				op:  token{Type: TokenIdentifier, Lexeme: "i"},
				exp: literalExpression{value: 0},
			},
			whileStatement{
				body: blockStmt{
					statements: []statement{
						printStmt{exp: variableExpression{op: token{Type: TokenIdentifier, Lexeme: "i"}}},
						expressionStmt{
							exp: assignExpression{
								exp: literalExpression{},
							},
						},
					},
				},
				cond: binaryExpression{
					left:  variableExpression{op: token{Type: TokenIdentifier, Lexeme: "i"}},
					op:    token{Type: TokenLess},
					right: literalExpression{value: 10},
				},
			},
		},
	}, stmt)
}
