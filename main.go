package main

func main() {

	env := newGlobalEnvironment()

	code := `

	print "Current time: " + now()

	num=10
	function fibonacci(n) {
		if n <= 1 {
			return n
		}
		return fibonacci(n-2) + fibonacci(n-1)
	}

	print "Fibonnacci of " + num + "=" + fibonacci(num)
	`
	sc := newScanner(code)
	tokens := sc.scanTokens()
	p := parser{
		tokens: tokens,
	}
	stmt := p.parse()
	for _, s := range stmt {
		s.evaluate(env)
	}

}
