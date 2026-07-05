package repl

import (
	"bufio"
	"fmt"

	"interpreter/evaluator"
	"interpreter/lexer"
	"interpreter/object"
	"interpreter/parser"
	"io"
)

const prompt = ">>"

const LOGO = `
     ,_,
    (o,o)
    {" '}
   -"---"-
`

func STart(in io.Reader, out io.Writer) {

	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Printf(prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}
func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, LOGO)
	io.WriteString(out, "woops! we ran into an error\n")
	io.WriteString(out, "parser error\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
