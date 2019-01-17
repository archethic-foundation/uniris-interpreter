package uniris

import "fmt"

//Interpret smart contract code
func Interpret(code string, env *Environment) (string, error) {

	globals := NewEnvironment(nil)
	globals.Set("now", currentTimestampFunc{})

	if env == nil {
		env = NewEnvironment(nil)
	}
	env.enclosing = globals

	sc := newScanner(code)
	tokens := sc.scanTokens()
	p := parser{
		tokens: tokens,
	}
	stmt, err := p.parse()
	if err != nil {
		return "", err
	}

	res := ""

	for _, s := range stmt {
		val, err := s.evaluate(env)
		if err != nil {
			return "", err
		}
		if val != nil {
			res += fmt.Sprintf("%v\n", val)
		}
	}

	return res, nil
}
