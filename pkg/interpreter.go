package uniris

import "fmt"

//Interpret smart contract code
func Interpret(code string) error {

	globals := &environment{}
	globals.set("now", currentTimestampFunc{})

	env := &environment{
		enclosing: globals,
	}

	sc := newScanner(code)
	tokens := sc.scanTokens()
	p := parser{
		tokens: tokens,
	}
	stmt, err := p.parse()
	if err != nil {
		return err
	}
	for _, s := range stmt {
		val, err := s.evaluate(env)
		if err != nil {
			return err
		}
		if val != nil {
			fmt.Println(val)
		}
	}

	return nil
}
