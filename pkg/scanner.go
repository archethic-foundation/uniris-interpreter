package uniris

import (
	"fmt"
	"strconv"
)

type token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

type TokenType string

var keywords = map[string]TokenType{
	"if":          TokenIf,
	"else":        TokenElse,
	"while":       TokenWhile,
	"for":         TokenFor,
	"or":          TokenOr,
	"and":         TokenAnd,
	"true":        TokenTrue,
	"false":       TokenFalse,
	"function":    TokenFunction,
	"print":       TokenPrint,
	"return":      TokenReturn,
	"transaction": TokenTransaction,
	"contract":    TokenContract,
}

const (
	// Single-character tokens.
	TokenLeftParenthesis  TokenType = "LEFT_PARENTHESIS"
	TokenRightParenthesis TokenType = "RIGHT_PARENTHESIS"
	TokenLeftBracket      TokenType = "LEFT_BRACKET"
	TokenRightBracket     TokenType = "RIGHT_BRACKET"
	TokenPlus             TokenType = "PLUS"
	TokenMinus            TokenType = "MINUS"
	TokenStar             TokenType = "STAR"
	TokenSlash            TokenType = "SLASH"
	TokenDot              TokenType = "DOT"
	TokenComma            TokenType = "COMMA"
	TokenSemiColon        TokenType = "SEMICOLON"

	//One or two character tokens
	TokenBang         TokenType = "BANG"
	TokenBangEqual    TokenType = "BANQ_EQUAL"
	TokenEqual        TokenType = "EQUAL"
	TokenEqualEqual   TokenType = "EQUAL_EQUAL"
	TokenLess         TokenType = "LESS"
	TokenGreater      TokenType = "GREATER"
	TokenLessEqual    TokenType = "LESS_EQUAL"
	TokenGreaterEqual TokenType = "GREATER_EQUAL"

	//Literals
	TokenIdentifier TokenType = "IDENTIFIER"
	TokenString     TokenType = "STRING"
	TokenNumber     TokenType = "NUMBER"

	//Keywords
	TokenPrint       TokenType = "PRINT"
	TokenIf          TokenType = "IF"
	TokenElse        TokenType = "ELSE"
	TokenAnd         TokenType = "AND"
	TokenOr          TokenType = "OR"
	TokenWhile       TokenType = "WHILE"
	TokenFor         TokenType = "FOR"
	TokenEndOfFile   TokenType = "EOF"
	TokenTrue        TokenType = "TRUE"
	TokenFalse       TokenType = "FALSE"
	TokenFunction    TokenType = "FUNC"
	TokenReturn      TokenType = "RETURN"
	TokenTransaction TokenType = "TRANSACTION"
	TokenContract    TokenType = "CONTRACT"
)

type scanner struct {
	source  []rune
	start   int
	current int
	line    int
	tokens  []token
}

func newScanner(code string) scanner {
	return scanner{
		source:  []rune(code),
		start:   0,
		current: 0,
		line:    1,
	}
}

func (sc *scanner) scanTokens() []token {
	for !sc.isAtEnd() {
		sc.start = sc.current
		sc.scanToken()
	}

	sc.tokens = append(sc.tokens, token{Type: TokenEndOfFile, Line: sc.line})
	return sc.tokens
}

func (sc *scanner) isAtEnd() bool {
	return sc.current >= len(sc.source)
}

func (sc *scanner) scanToken() {
	c := sc.advance()
	switch c {
	case '(':
		sc.addEmptyToken(TokenLeftParenthesis)
		break
	case ')':
		sc.addEmptyToken(TokenRightParenthesis)
		break
	case '{':
		sc.addEmptyToken(TokenLeftBracket)
		break
	case '}':
		sc.addEmptyToken(TokenRightBracket)
		break
	case '+':
		sc.addEmptyToken(TokenPlus)
		break
	case '-':
		sc.addEmptyToken(TokenMinus)
		break
	case '*':
		sc.addEmptyToken(TokenStar)
		break
	case '.':
		sc.addEmptyToken(TokenDot)
		break
	case ',':
		sc.addEmptyToken(TokenComma)
		break
	case ';':
		sc.addEmptyToken(TokenSemiColon)
		break
	case '!':
		if sc.match('=') {
			sc.addEmptyToken(TokenBangEqual)
		} else {
			sc.addEmptyToken(TokenBang)
		}
		break
	case '/':
		if sc.match('/') {
			// A comment goes until the end of the line.
			for sc.peek() != '\n' && !sc.isAtEnd() {
				sc.advance()
			}
			break
		}
		sc.addEmptyToken(TokenSlash)
		break
	case '=':
		if sc.match('=') {
			sc.addEmptyToken(TokenEqualEqual)
		} else {
			sc.addEmptyToken(TokenEqual)
		}
		break
	case '>':
		if sc.match('=') {
			sc.addEmptyToken(TokenGreaterEqual)
		} else {
			sc.addEmptyToken(TokenGreater)
		}
		break
	case '<':
		if sc.match('=') {
			sc.addEmptyToken(TokenLessEqual)
		} else {
			sc.addEmptyToken(TokenLess)
		}
		break
	case ' ':
	case '\r':
	case '\t':
		// Ignore whitespace.
		break
	case '\n':
		sc.line++
		break
	case '"':
		sc.string()
		break
	default:
		if sc.isDigit(c) {
			sc.number()
		} else if sc.isAlpha(c) {
			sc.identifier()
		} else {
			panic(fmt.Errorf("Error: Line: %d, Unexpected character: %s", sc.line, string(c)))
		}
		break
	}
}

func (sc *scanner) isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (sc *scanner) isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func (sc *scanner) isAlphaNumeric(c rune) bool {
	return sc.isAlpha(c) || sc.isDigit(c)
}

func (sc *scanner) identifier() {
	for sc.isAlphaNumeric(sc.peek()) {
		sc.advance()
	}

	text := sc.source[sc.start:sc.current]
	if tokenType, exist := keywords[string(text)]; exist {
		sc.addEmptyToken(tokenType)
	} else {
		sc.addEmptyToken(TokenIdentifier)
	}
}

func (sc *scanner) number() {
	for sc.isDigit(sc.peek()) {
		sc.advance()
	}

	// Look for a fractional part.
	if sc.peek() == '.' && sc.isDigit(sc.peekNext()) {
		// Consume the "."
		sc.advance()

		for sc.isDigit(sc.peek()) {
			sc.advance()
		}
	}

	float, err := strconv.ParseFloat(string(sc.source[sc.start:sc.current]), 64)
	if err == nil {
		sc.addToken(TokenNumber, float)
	}
}

func (sc *scanner) peek() rune {
	if sc.isAtEnd() {
		return rune(-1)
	}
	return sc.source[sc.current]
}

func (sc *scanner) peekNext() rune {
	if sc.current+1 >= len(sc.source) {
		return rune(-1)
	}
	return sc.source[sc.current+1]
}

func (sc *scanner) match(c rune) bool {
	if sc.isAtEnd() {
		return false
	}
	if sc.source[sc.current] != c {
		return false
	}
	sc.current++
	return true
}

func (sc *scanner) addToken(t TokenType, lit interface{}) {
	text := sc.source[sc.start:sc.current]
	sc.tokens = append(sc.tokens, token{
		Type:    t,
		Lexeme:  string(text),
		Literal: lit,
		Line:    sc.line,
	})
}

func (sc *scanner) addEmptyToken(t TokenType) {
	sc.addToken(t, nil)
}

func (sc *scanner) advance() rune {
	c := sc.source[sc.current]
	sc.current++
	return c
}

func (sc *scanner) string() {

	for sc.peek() != '"' && !sc.isAtEnd() {
		if sc.peek() == '\n' {
			sc.line++
		}
		sc.advance()

	}

	// Unterminated string.
	if sc.isAtEnd() {
		panic(fmt.Sprintf("ERROR: Line: %d, Unterminated string.", sc.line))
	}

	// The closing ".
	sc.advance()

	// Trim the surrounding quotes.
	value := sc.source[sc.start+1 : sc.current-1]
	sc.addToken(TokenString, string(value))
}
