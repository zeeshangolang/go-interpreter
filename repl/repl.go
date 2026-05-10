package repl

import (
	"bufio"
	"fmt"
	"interpreter/lexer"
	"interpreter/tokens"
	"io"
)

const prompt = ">>"

func STart(in io.Reader, ot io.Writer) {

	scanner := bufio.NewScanner(in)

	scanned := scanner.Scan()
	if !scanned {
		return
	}

	line := scanner.Text()
	l := lexer.New(line)

	for tok := l.NextToken(); tok.Type != tokens.EOF; tok = l.NextToken() {
		fmt.Fprintf(ot, "%+v\n", tok)
	}
}
