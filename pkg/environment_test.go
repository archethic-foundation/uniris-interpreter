package uniris

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetEnvValue(t *testing.T) {
	e := NewEnvironment(nil)
	e.set("a", 2)
	assert.Equal(t, 2, e.values["a"])
}

func TestSetEnclosingEnvValue(t *testing.T) {

	enc := NewEnvironment(nil)
	enc.set("a", 2)

	e := NewEnvironment(enc)

	assert.Nil(t, e.values["a"])
	assert.Equal(t, 2, e.enclosing.values["a"])

	e.set("a", 5)
	assert.Equal(t, 5, e.enclosing.values["a"])

	e.set("b", 10)
	assert.Equal(t, 10, e.values["b"])
	assert.Nil(t, e.enclosing.values["b"])
}

func TestGetEnvValue(t *testing.T) {
	e := NewEnvironment(nil)
	e.set("a", 2)
	val, err := e.get("a")
	assert.Nil(t, err)
	assert.Equal(t, 2, val)
}

func TestGetEnclosingEnvValue(t *testing.T) {
	enc := NewEnvironment(nil)
	enc.set("a", 2)

	e := NewEnvironment(enc)

	val, err := e.get("a")
	assert.Nil(t, err)
	assert.Equal(t, 2, val)
}

func TestGetUndefinedEnvValue(t *testing.T) {
	e := NewEnvironment(nil)
	_, err := e.get("a")
	assert.Error(t, err, "Undefined variable a")
}
