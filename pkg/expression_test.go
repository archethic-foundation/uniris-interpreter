package uniris

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLiteralExpression(t *testing.T) {
	e := literalExpression{
		value: 10,
	}
	val, err := e.evaluate(NewEnvironment(nil))
	assert.Nil(t, err)
	assert.Equal(t, 10, val)
}

func TestAssignExpression(t *testing.T) {
	e := assignExpression{
		exp: literalExpression{
			value: 10,
		},
		op: token{
			Lexeme: "a",
		},
	}

	env := NewEnvironment(nil)
	_, err := e.evaluate(env)
	assert.Nil(t, err)
	val, err := env.get("a")
	assert.Nil(t, err)
	assert.Equal(t, 10, val)
}

func TestVariableExpression(t *testing.T) {

	env := NewEnvironment(nil)
	env.set("a", 10)

	e := variableExpression{
		op: token{
			Lexeme: "a",
		},
	}

	val, err := e.evaluate(env)
	assert.Nil(t, err)
	assert.Equal(t, 10, val)
}

func TestBinaryMinusExpression(t *testing.T) {
	e := binaryExpression{
		left: literalExpression{
			value: float64(10),
		},
		right: literalExpression{
			value: float64(10),
		},
		op: token{Type: TokenMinus},
	}

	val, err := e.evaluate(NewEnvironment(nil))
	assert.Nil(t, err)
	assert.Equal(t, float64(0), val)
}

func TestBinaryPlusExpression(t *testing.T) {
	e := binaryExpression{
		left: literalExpression{
			value: float64(10),
		},
		right: literalExpression{
			value: float64(10),
		},
		op: token{Type: TokenPlus},
	}

	val, err := e.evaluate(NewEnvironment(nil))
	assert.Nil(t, err)
	assert.Equal(t, float64(20), val)

	e = binaryExpression{
		left: literalExpression{
			value: "hello ",
		},
		right: literalExpression{
			value: "world",
		},
		op: token{Type: TokenPlus},
	}
	val, err = e.evaluate(NewEnvironment(nil))
	assert.Nil(t, err)
	assert.Equal(t, "hello world", val)
}

func TestBinaryStarExpression(t *testing.T) {
	e := binaryExpression{
		left: literalExpression{
			value: float64(10),
		},
		right: literalExpression{
			value: float64(10),
		},
		op: token{Type: TokenStar},
	}

	val, err := e.evaluate(NewEnvironment(nil))
	assert.Nil(t, err)
	assert.Equal(t, float64(100), val)
}

func TestBinarySlashExpression(t *testing.T) {
	e := binaryExpression{
		left: literalExpression{
			value: float64(10),
		},
		right: literalExpression{
			value: float64(10),
		},
		op: token{Type: TokenSlash},
	}

	val, err := e.evaluate(NewEnvironment(nil))
	assert.Nil(t, err)
	assert.Equal(t, float64(1), val)
}

func TestBinaryGreaterExpression(t *testing.T) {
	e := binaryExpression{
		left: literalExpression{
			value: float64(11),
		},
		right: literalExpression{
			value: float64(10),
		},
		op: token{Type: TokenGreater},
	}

	val, err := e.evaluate(NewEnvironment(nil))
	assert.Nil(t, err)
	assert.Equal(t, true, val)
}

func TestBinaryGreaterEqualExpression(t *testing.T) {
	e := binaryExpression{
		left: literalExpression{
			value: float64(10),
		},
		right: literalExpression{
			value: float64(10),
		},
		op: token{Type: TokenGreaterEqual},
	}

	val, err := e.evaluate(NewEnvironment(nil))
	assert.Nil(t, err)
	assert.Equal(t, true, val)
}

func TestBinaryLessExpression(t *testing.T) {
	e := binaryExpression{
		left: literalExpression{
			value: float64(10),
		},
		right: literalExpression{
			value: float64(10),
		},
		op: token{Type: TokenLess},
	}

	val, err := e.evaluate(NewEnvironment(nil))
	assert.Nil(t, err)
	assert.Equal(t, false, val)
}

func TestBinaryLessEqualExpression(t *testing.T) {
	e := binaryExpression{
		left: literalExpression{
			value: float64(10),
		},
		right: literalExpression{
			value: float64(10),
		},
		op: token{Type: TokenLessEqual},
	}

	val, err := e.evaluate(NewEnvironment(nil))
	assert.Nil(t, err)
	assert.Equal(t, true, val)
}

func TestBinaryEqualEqualExpression(t *testing.T) {
	e := binaryExpression{
		left: literalExpression{
			value: float64(10),
		},
		right: literalExpression{
			value: float64(10),
		},
		op: token{Type: TokenEqualEqual},
	}

	val, err := e.evaluate(NewEnvironment(nil))
	assert.Nil(t, err)
	assert.Equal(t, true, val)
}

func TestBinaryBangEqualExpression(t *testing.T) {
	e := binaryExpression{
		left: literalExpression{
			value: float64(10),
		},
		right: literalExpression{
			value: float64(10),
		},
		op: token{Type: TokenBangEqual},
	}

	val, err := e.evaluate(NewEnvironment(nil))
	assert.Nil(t, err)
	assert.Equal(t, false, val)
}

func TestGroupExpression(t *testing.T) {
	e := groupingExpression{
		exp: literalExpression{
			value: 10,
		},
	}

	val, err := e.evaluate(NewEnvironment(nil))
	assert.Nil(t, err)
	assert.Equal(t, 10, val)
}

func TestUnaryBangExpression(t *testing.T) {
	e := unaryExpression{
		op: token{
			Type: TokenBang,
		},
		right: literalExpression{
			value: true,
		},
	}

	val, err := e.evaluate(NewEnvironment(nil))
	assert.Nil(t, err)
	assert.Equal(t, false, val)

}

func TestUnaryMinusExpression(t *testing.T) {
	e := unaryExpression{
		op: token{
			Type: TokenMinus,
		},
		right: literalExpression{
			value: float64(10),
		},
	}

	val, err := e.evaluate(NewEnvironment(nil))
	assert.Nil(t, err)
	assert.Equal(t, float64(-10), val)
}

func TestLogicalANDExpression(t *testing.T) {
	e := logicalExpression{
		left: literalExpression{
			value: true,
		},
		right: literalExpression{
			value: 10,
		},
		op: token{
			Type: TokenAnd,
		},
	}

	val, err := e.evaluate(NewEnvironment(nil))
	assert.Nil(t, err)
	assert.Equal(t, 10, val)

	e = logicalExpression{
		left: literalExpression{
			value: false,
		},
		right: literalExpression{
			value: 10,
		},
		op: token{
			Type: TokenAnd,
		},
	}

	val, err = e.evaluate(NewEnvironment(nil))
	assert.Nil(t, err)
	assert.Equal(t, false, val)
}

func TestLogicalORExpression(t *testing.T) {
	e := logicalExpression{
		left: literalExpression{
			value: true,
		},
		right: literalExpression{
			value: 10,
		},
		op: token{
			Type: TokenOr,
		},
	}

	val, err := e.evaluate(NewEnvironment(nil))
	assert.Nil(t, err)
	assert.Equal(t, true, val)

	e = logicalExpression{
		left: literalExpression{
			value: false,
		},
		right: literalExpression{
			value: 10,
		},
		op: token{
			Type: TokenOr,
		},
	}

	val, err = e.evaluate(NewEnvironment(nil))
	assert.Nil(t, err)
	assert.Equal(t, 10, val)
}

func TestCallExpression(t *testing.T) {
	e := callExpression{
		args: []expression{
			literalExpression{
				value: 10,
			},
		},
		callee: testFuncExpression{},
	}

	val, err := e.evaluate(NewEnvironment(nil))
	assert.Nil(t, err)
	assert.Equal(t, 10, val)
}

func TestCallNotCallableExpression(t *testing.T) {
	e := callExpression{
		callee: literalExpression{},
	}

	_, err := e.evaluate(NewEnvironment(nil))
	assert.Error(t, err, "Can only call functions")
}

type testFuncExpression struct {
}

func (f testFuncExpression) evaluate(env *Environment) (interface{}, error) {
	return testFuncCallable{}, nil
}

type testFuncCallable struct{}

func (f testFuncCallable) call(env *Environment, args ...interface{}) (res interface{}, err error) {
	return args[0], nil
}
