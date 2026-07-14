package parser

import (
	"fmt"
	"interpreter/ast"
	"interpreter/lexer"
	"interpreter/tokens"
	"strconv"
	//"monkey/ast"
)

var precedenes = map[tokens.TokenType]int{

	tokens.EQ:       EQUALS,
	tokens.NEQ:      EQUALS,
	tokens.LTHAN:    LESSGREATER,
	tokens.GTHAN:    LESSGREATER,
	tokens.PLUS:     SUM,
	tokens.MINUS:    SUM,
	tokens.SLASH:    PRODUCT,
	tokens.ASTERISK: PRODUCT,
	tokens.LPARAN:   CALL,
	tokens.LBRACKET: INDEX,
}

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       //array[index]
)

type Parser struct {
	l         *lexer.Lexer
	errors    []string
	currToken tokens.Token
	peekToken tokens.Token

	prefixParserfn map[tokens.TokenType]prefixParserfn
	infixParserfn  map[tokens.TokenType]infixParserfn
}

type (
	prefixParserfn func() ast.Expression
	infixParserfn  func(ast.Expression) ast.Expression
)

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{}}

	p.prefixParserfn = make(map[tokens.TokenType]prefixParserfn)
	p.infixParserfn = make(map[tokens.TokenType]infixParserfn)
	p.RegisterPrefix(tokens.STRING, p.ParseStringLiteral)
	p.RegisterPrefix(tokens.LBRACE, p.parseHashLiterals)
	p.RegisterPrefix(tokens.LBRACKET, p.parseArrayLiteral)
	p.RegisterPrefix(tokens.FUNCTION, p.parseFunctionLiteral)
	p.RegisterInfix(tokens.LPARAN, p.parseCallExpression)
	p.RegisterPrefix(tokens.LPARAN, p.parseGroupedExpression)
	p.RegisterPrefix(tokens.TRUE, p.parseBoolean)
	p.RegisterPrefix(tokens.FALSE, p.parseBoolean)
	p.RegisterPrefix(tokens.IF, p.parseIfExpression)
	p.RegisterInfix(tokens.PLUS, p.parseInfixExpression)
	p.RegisterInfix(tokens.LBRACKET, p.parseIndexExpression)
	p.RegisterInfix(tokens.MINUS, p.parseInfixExpression)
	p.RegisterInfix(tokens.SLASH, p.parseInfixExpression)
	p.RegisterInfix(tokens.ASTERISK, p.parseInfixExpression)
	p.RegisterInfix(tokens.EQ, p.parseInfixExpression)
	p.RegisterInfix(tokens.NEQ, p.parseInfixExpression)
	p.RegisterInfix(tokens.LTHAN, p.parseInfixExpression)
	p.RegisterInfix(tokens.GTHAN, p.parseInfixExpression)
	p.RegisterPrefix(tokens.BANG, p.parsePrefixExpression)
	p.RegisterPrefix(tokens.MINUS, p.parsePrefixExpression)
	p.RegisterPrefix(tokens.IDENT, p.parseIdentifiers)
	p.RegisterPrefix(tokens.INT, p.parseIntegerLiteral)

	p.NextToken()
	p.NextToken()
	return p
}

func (p *Parser) parseHashLiterals() ast.Expression {
	hash := &ast.HashLiterals{Token: p.currToken}
	hash.Piars = make(map[ast.Expression]ast.Expression)

	for !p.peekTOkenIs(tokens.RBARCE) {
		p.NextToken()
		key := p.parseExpression(LOWEST)
		if !p.expectedPeek(tokens.COLON) {
			return nil
		}
		p.NextToken()
		value := p.parseExpression(LOWEST)

		hash.Piars[key] = value

		if !p.peekTOkenIs(tokens.RBARCE) && !p.expectedPeek(tokens.COMMA) {
			return nil
		}
	}
	if !p.expectedPeek(tokens.RBARCE) {
		return nil
	}
	return hash
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {

	exp := &ast.IndexExpression{Token: p.currToken, Left: left}

	p.NextToken()

	exp.Index = p.parseExpression(LOWEST)

	if !p.expectedPeek(tokens.RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) parseArrayLiteral() ast.Expression {

	array := &ast.ArrayLiteral{Token: p.currToken}

	array.Elements = p.parseExpressionList(tokens.RBRACKET)

	return array
}

func (p *Parser) parseExpressionList(end tokens.TokenType) []ast.Expression {

	list := []ast.Expression{}

	if p.peekTOkenIs(end) {
		p.NextToken()
		return list
	}

	p.NextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTOkenIs(tokens.COMMA) {
		p.NextToken()
		p.NextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectedPeek(end) {
		return nil
	}

	return list

}
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.currToken, Function: function}
	exp.Arguments = p.parseExpressionList(tokens.RPARAN)
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTOkenIs(tokens.RPARAN) {
		p.NextToken()
		return args
	}

	p.NextToken()

	args = append(args, p.parseExpression(LOWEST))

	for p.peekTOkenIs(tokens.COMMA) {
		p.NextToken()
		p.NextToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	if !p.expectedPeek(tokens.RPARAN) {
		return nil
	}

	return args
}

func (p *Parser) ParseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.currToken}

	if !p.expectedPeek(tokens.LPARAN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectedPeek(tokens.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatment()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {

	identifiers := []*ast.Identifier{}

	if p.peekTOkenIs(tokens.RPARAN) {
		p.NextToken()
		return identifiers
	}

	p.NextToken()

	ident := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTOkenIs(tokens.COMMA) {
		p.NextToken()
		p.NextToken()
		ident := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
		identifiers = append(identifiers, ident)
	}

	for !p.expectedPeek(tokens.RPARAN) {
		return nil
	}

	return identifiers

}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.currToken}

	if !p.expectedPeek(tokens.LPARAN) {
		return nil
	}

	p.NextToken()

	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectedPeek(tokens.RPARAN) {
		return nil
	}

	if !p.expectedPeek(tokens.LBRACE) {
		return nil
	}

	expression.Consequences = p.parseBlockStatment()

	if p.peekTOkenIs(tokens.ELSE) {
		p.NextToken()
		if !p.expectedPeek(tokens.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatment()
	}

	return expression
}

func (p *Parser) parseBlockStatment() *ast.BlockStatement {

	block := &ast.BlockStatement{Token: p.currToken}
	block.Statement = []ast.Statement{}

	p.NextToken()

	for !p.currTokenIs(tokens.RBARCE) && !p.currTokenIs(tokens.EOF) {
		stmt := p.parseStament()
		if stmt != nil {
			block.Statement = append(block.Statement, stmt)
		}
		p.NextToken()
	}
	return block

}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.NextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectedPeek(tokens.RPARAN) {
		return nil

	}

	return exp
}

func (p *Parser) parseBoolean() ast.Expression {

	return &ast.Boolean{Token: p.currToken, Value: p.currTokenIs(tokens.TRUE)}
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {

	expression := &ast.InfixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
		Left:     left,
	}

	precedence := p.currPrecedence()
	p.NextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
	}

	p.NextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseIdentifiers() ast.Expression {

	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t tokens.TokenType) {
	msg := fmt.Sprintf("expected nex token to be %s got: %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) NextToken() {

	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()

}

func (p *Parser) RegisterPrefix(TokenType tokens.TokenType, fn prefixParserfn) {

	p.prefixParserfn[TokenType] = fn
}

func (p *Parser) RegisterInfix(TokenType tokens.TokenType, fn infixParserfn) {
	p.infixParserfn[TokenType] = fn
}

func (p *Parser) ParseProgram() *ast.Program {

	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currToken.Type != tokens.EOF {
		stmt := p.parseStament()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.NextToken()
	}

	return program

}

func (p *Parser) parseStament() ast.Statement {
	switch p.currToken.Type {
	case tokens.LET:
		return p.parseLetStatement()
	case tokens.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}

}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currToken}

	if !p.expectedPeek(tokens.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	if !p.expectedPeek(tokens.ASSIGN) {
		return nil
	}

	p.NextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTOkenIs(tokens.SEMICOLON) {
		p.NextToken()
	}

	return stmt

}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {

	stmt := &ast.ReturnStatement{Token: p.currToken}

	p.NextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	for !p.currTokenIs(tokens.SEMICOLON) {
		p.NextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTOkenIs(tokens.SEMICOLON) {
		p.NextToken()
	}
	return stmt
}

func (p *Parser) currTokenIs(token tokens.TokenType) bool {
	return p.currToken.Type == token
}

func (p *Parser) peekTOkenIs(t tokens.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectedPeek(token tokens.TokenType) bool {
	if p.peekTOkenIs(token) {
		p.NextToken()
		return true
	} else {
		p.peekError(token)
		return false
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {

	lit := &ast.IntegerLiteral{Token: p.currToken}

	value, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit

}

func (p *Parser) noPrefixParseFnError(t tokens.TokenType) {

	msg := fmt.Sprintf("no prefix parse function for the %s Found", t)

	p.errors = append(p.errors, msg)

}

func (p *Parser) peekPrecedence() int {
	if precedence, ok := precedenes[p.peekToken.Type]; ok {
		return precedence
	}

	return LOWEST
}

func (p *Parser) currPrecedence() int {
	if precedence, ok := precedenes[p.currToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

func (p *Parser) parseExpression(precedene int) ast.Expression {

	prefix := p.prefixParserfn[p.currToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.currToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.peekTOkenIs(tokens.SEMICOLON) && precedene < p.peekPrecedence() {
		infix := p.infixParserfn[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.NextToken()

		leftExp = infix(leftExp)

	}

	return leftExp
}
