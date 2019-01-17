package uniris

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitializeScanner(t *testing.T) {
	s := newScanner("print")
	assert.Equal(t, 0, s.current)
	assert.Equal(t, 1, s.line)
	assert.Len(t, s.source, 5)
}

func TestScanAdvance(t *testing.T) {
	s := newScanner("hello")
	s.advance()
	assert.Equal(t, 1, s.current)
	s.advance()
	assert.Equal(t, 2, s.current)
}

func TestIsAtEnd(t *testing.T) {
	s := newScanner("")
	s.current = 1
	assert.True(t, s.isAtEnd())

	s = newScanner("print")
	assert.False(t, s.isAtEnd())
}

func TestAddEmptyToken(t *testing.T) {
	s := newScanner("(")
	s.addEmptyToken(TokenLeftParenthesis)
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, TokenLeftParenthesis, s.tokens[0].Type)
}

func TestScanIsAlpha(t *testing.T) {
	s := newScanner("")
	assert.True(t, s.isAlpha(rune('a')))
	assert.False(t, s.isAlpha(rune(10)))
}

func TestScanIsDigit(t *testing.T) {
	s := newScanner("")
	assert.True(t, s.isDigit(rune('1')))
}

func TestScanIsAlphanumeric(t *testing.T) {
	s := newScanner("")
	assert.True(t, s.isAlphaNumeric(rune('a')))
	assert.True(t, s.isAlphaNumeric(rune('1')))
}

func TestScanIdentifier(t *testing.T) {
	s := newScanner("a")
	s.identifier()
	assert.Equal(t, TokenIdentifier, s.tokens[0].Type)

	s = newScanner("if")
	s.identifier()
	assert.Equal(t, TokenIf, s.tokens[0].Type)
}

func TestScanString(t *testing.T) {
	s := newScanner("\"hello\"")
	s.advance()
	s.string()
	assert.Equal(t, TokenString, s.tokens[0].Type)
	assert.Equal(t, "hello", s.tokens[0].Literal)
	assert.Equal(t, "\"hello\"", s.tokens[0].Lexeme)

	s = newScanner("\"hello")
	s.advance()
	assert.Panics(t, s.string)

	s = newScanner("\"hello\nWorld\"")
	s.advance()
	s.string()
	assert.Equal(t, 2, s.line)
}

func TestScanNumber(t *testing.T) {
	s := newScanner("123")
	s.number()
	assert.Equal(t, "123", s.tokens[0].Lexeme)
	assert.Equal(t, TokenNumber, s.tokens[0].Type)
	assert.Equal(t, float64(123), s.tokens[0].Literal)
}

func TestScanTokenParenthesis(t *testing.T) {
	s := newScanner("(")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, TokenLeftParenthesis, s.tokens[0].Type)
	assert.Equal(t, 1, s.tokens[0].Line)

	s = newScanner(")")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, TokenRightParenthesis, s.tokens[0].Type)

}

func TestScanTokenBrackets(t *testing.T) {
	s := newScanner("{")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, TokenLeftBracket, s.tokens[0].Type)

	s = newScanner("}")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, TokenRightBracket, s.tokens[0].Type)
}

func TestScanTokenArithmetic(t *testing.T) {
	s := newScanner("+")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, TokenPlus, s.tokens[0].Type)

	s = newScanner("-")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, TokenMinus, s.tokens[0].Type)

	s = newScanner("/")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, TokenSlash, s.tokens[0].Type)

	s = newScanner("*")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, TokenStar, s.tokens[0].Type)
}

func TestScanTokenBang(t *testing.T) {
	s := newScanner("!")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, TokenBang, s.tokens[0].Type)

	s = newScanner("!=")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, TokenBangEqual, s.tokens[0].Type)

}

func TestScanTokenEqual(t *testing.T) {
	s := newScanner("=")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, TokenEqual, s.tokens[0].Type)

	s = newScanner("==")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, TokenEqualEqual, s.tokens[0].Type)
}

func TestScanTokenComparison(t *testing.T) {
	s := newScanner(">")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, TokenGreater, s.tokens[0].Type)

	s = newScanner(">=")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, TokenGreaterEqual, s.tokens[0].Type)

	s = newScanner("<")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, TokenLess, s.tokens[0].Type)

	s = newScanner("<=")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, TokenLessEqual, s.tokens[0].Type)
}

func TestScanTokenNumber(t *testing.T) {
	s := newScanner("1")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, TokenNumber, s.tokens[0].Type)

	s = newScanner("1.10")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, TokenNumber, s.tokens[0].Type)
}

func TestScanTokenIdentifier(t *testing.T) {
	s := newScanner("a")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, TokenIdentifier, s.tokens[0].Type)
}

func TestScanTokenString(t *testing.T) {
	s := newScanner("\"hello\"")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, TokenString, s.tokens[0].Type)
}

func TestScanComments(t *testing.T) {
	s := newScanner("//a")
	s.scanToken()
	assert.Len(t, s.tokens, 0)
}

func TestScanTokenDot(t *testing.T) {
	s := newScanner(".")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, TokenDot, s.tokens[0].Type)

}

func TestScanTokenComma(t *testing.T) {
	s := newScanner(",")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, TokenComma, s.tokens[0].Type)
}

func TestScanTokenSemiColon(t *testing.T) {
	s := newScanner(";")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, TokenSemiColon, s.tokens[0].Type)
}

func TestScanTokenWhitespace(t *testing.T) {
	s := newScanner(" ")
	s.scanToken()
	assert.Len(t, s.tokens, 0)

	s = newScanner("\r")
	s.scanToken()
	assert.Len(t, s.tokens, 0)

	s = newScanner("\t")
	s.scanToken()
	assert.Len(t, s.tokens, 0)
}

func TestScanTokenBreakline(t *testing.T) {
	s := newScanner("\n")
	s.scanToken()
	assert.Len(t, s.tokens, 0)
	assert.Equal(t, 2, s.line)
}

func TestScanTokenUnexpected(t *testing.T) {
	s := newScanner("°")
	assert.Panics(t, s.scanToken)
}

func TestScanMultipleTokens(t *testing.T) {
	s := newScanner("print 2+2")
	tokens := s.scanTokens()
	assert.Len(t, tokens, 5)
	assert.Equal(t, TokenPrint, tokens[0].Type)
	assert.Equal(t, TokenNumber, tokens[1].Type)
	assert.Equal(t, float64(2), tokens[1].Literal)
	assert.Equal(t, TokenPlus, tokens[2].Type)
	assert.Equal(t, TokenNumber, tokens[3].Type)
	assert.Equal(t, float64(2), tokens[3].Literal)
	assert.Equal(t, TokenEndOfFile, tokens[4].Type)
}
