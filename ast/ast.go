package ast

import "interpreter/tokens"

type Node interface {
	Tokenliteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) Tokenliteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].Tokenliteral()
	} else {
		return ""
	}
}

type LetStatement struct {
	Token tokens.Token
	Name  *Identifier
	Value Expression
}
