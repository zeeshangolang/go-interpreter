package lexer

import (
	"interpreter/tokens"
	//	"main/tokenss"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
    let ten  = 10;

    let add = fn(x, y) {
    x + y
    };
    let result = add(five, ten);
	!-/*5
	5 < 10 > 5

	if (5 < 10) {
	return true;
	} else {
	 return false
	 }

	 10 == 10
	 10 != 10

    `

	testss := []struct {
		expectedType    tokens.TokenType
		expectedLiteral string
	}{
		{tokens.LET, "let"},
		{tokens.IDENT, "five"},
		{tokens.ASSIGN, "="},
		{tokens.INT, "5"},
		{tokens.SEMICOLON, ";"},
		{tokens.LET, "let"},
		{tokens.IDENT, "ten"},
		{tokens.ASSIGN, "="},
		{tokens.INT, "10"},
		{tokens.SEMICOLON, ";"},
		{tokens.LET, "let"},
		{tokens.IDENT, "add"},
		{tokens.ASSIGN, "="},
		{tokens.FUNCTION, "fn"},
		{tokens.LPARAN, "("},
		{tokens.IDENT, "x"},
		{tokens.COMMA, ","},
		{tokens.IDENT, "y"},
		{tokens.RPARAN, ")"},
		{tokens.LBRACE, "{"},
		{tokens.IDENT, "x"},
		{tokens.PLUS, "+"},
		{tokens.IDENT, "y"},
		//{tokens.SEMICOLON, ";"},
		{tokens.RBARCE, "}"},
		{tokens.SEMICOLON, ";"},
		{tokens.LET, "let"},
		{tokens.IDENT, "result"},
		{tokens.ASSIGN, "="},
		{tokens.IDENT, "add"},
		{tokens.LPARAN, "("},
		{tokens.IDENT, "five"},
		{tokens.COMMA, ","},
		{tokens.IDENT, "ten"},
		{tokens.RPARAN, ")"},
		{tokens.SEMICOLON, ";"},
		{tokens.BANG, "!"},
		{tokens.MINUS, "-"},
		{tokens.SLASH, "/"},
		{tokens.ASTERISK, "*"},
		{tokens.INT, "5"},
		{tokens.INT, "5"},
		{tokens.LTHAN, "<"},
		{tokens.INT, "10"},
		{tokens.GTHAN, ">"},
		{tokens.INT, "5"},
		{tokens.IF, "if"},
		{tokens.LPARAN, "("},
		{tokens.INT, "5"},
		{tokens.LTHAN, "<"},
		{tokens.INT, "10"},
		{tokens.RPARAN, ")"},
		{tokens.LBRACE, "{"},
		{tokens.RETURN, "return"},
		{tokens.TRUE, "true"},
		{tokens.SEMICOLON, ";"},
		{tokens.RBARCE, "}"},
		{tokens.ELSE, "else"},
		{tokens.LBRACE, "{"},
		{tokens.RETURN, "return"},
		{tokens.FALSE, "false"},
		{tokens.RBARCE, "}"},
		{tokens.INT, "10"},
		{tokens.EQ, "=="},
		{tokens.INT, "10"},
		{tokens.INT, "10"},
		{tokens.NEQ, "!="},
		{tokens.INT, "10"},
		{tokens.EOF, ""},
	}
	//l  := New(input)
	l := New(input)

	for i, tt := range testss {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("test[%d], expexted=%q <> got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] expected %q <> got %q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
