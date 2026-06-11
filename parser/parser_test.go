package parser

import (
	"fmt"
	"interpreter/ast"
	"interpreter/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
    let x = 5;
    let y= 10;
    let zee = 13;`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	CheckParserError(t, p)
	if program == nil {
		t.Fatalf("Parsed program returend nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program doesnt have enoug states bra , got:%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"zee"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}

}

func CheckParserError(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}
	t.Errorf("parser had %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("error ppa %q", msg)
	}
	t.FailNow()
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not let got:[%q]", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s ast not let statement got %q", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letstmt.name.value not good expected %q, got %q", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("name.TOkenLiteral not %q, got %q", name, letStmt.TokenLiteral())
		return false
	}

	return true

}

func TestReturnStatement(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 134123;`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	CheckParserError(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("prgrm doesnt contain 3 got %d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnstatement, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("stmt not *ast.Returnstat got %q", stmt)
			continue
		}

		if returnstatement.TokenLiteral() != "return" {

			t.Errorf("return statement tokenliteral not return got %q", returnstatement.TokenLiteral())
		}
	}

}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	CheckParserError(t, p)

	if len(program.Statements) != 1 {

		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program statement 0 not an ast.expression got %T", program.Statements[0])

	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Errorf("stmt.expression not %s. got=%s", "foobar", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s got %s", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	CheckParserError(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5", literal.TokenLiteral())
	}
}

func TestParsinPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		CheckParserError(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("progrm stmts does not contain %d stmts. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s. got %s",
				tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
			//60, 62
		}

	}

}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {

	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.integerliteral Got =%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d Got %d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.Literal Not %d Got %s", value, integ.TokenLiteral())
		return false
	}

	return true
}

func TestParsingInfixExpression(t *testing.T) {

	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		CheckParserError(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("programming.Satements does not contain %d statements got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])

		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfxExpression, got =%T", stmt.Expression)
		}

		if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %s. got =%s",
				tt.operator, exp.Operator)

		}

		if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}
}
