package uniris

//Interpret smart contract code
func Interpret(code string, env *environment) error {

	globals := &environment{}
	globals.set("now", currentTimestampFunc{})

	if env == nil {
		env = &environment{
			enclosing: globals,
		}
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
		_, err := s.evaluate(env)
		if err != nil {
			return err
		}
	}

	return nil
}
