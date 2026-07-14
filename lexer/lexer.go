package lexer

import (
	"interpreter/tokens"
	//"main/tokens"
)

type Lexer struct {
	Input        string
	Position     int
	Readposition int
	ch           byte
}

func (l *Lexer) NextToken() tokens.Token {
	var tok tokens.Token

	l.SkipWhiteSpaces()
	switch l.ch {

	case '-':
		tok = newToken(tokens.MINUS, l.ch)

	case '/':
		tok = newToken(tokens.SLASH, l.ch)
	case '*':
		tok = newToken(tokens.ASTERISK, l.ch)
	case ':':
		tok = newToken(tokens.COLON, l.ch)
	case '<':
		tok = newToken(tokens.LTHAN, l.ch)
	case '>':
		tok = newToken(tokens.GTHAN, l.ch)
	case ',':
		tok = newToken(tokens.COMMA, l.ch)
	case '+':
		tok = newToken(tokens.PLUS, l.ch)
	case '"':
		tok.Type = tokens.STRING
		tok.Literal = l.readString()
	case '(':
		tok = newToken(tokens.LPARAN, l.ch)
	case ')':
		tok = newToken(tokens.RPARAN, l.ch)
	case '{':
		tok = newToken(tokens.LBRACE, l.ch)
	case '}':
		tok = newToken(tokens.RBARCE, l.ch)
	case '[':
		tok = newToken(tokens.LBRACKET, l.ch)
	case ']':
		tok = newToken(tokens.RBRACKET, l.ch)
	case ';':
		tok = newToken(tokens.SEMICOLON, l.ch)

	case 0:
		tok.Literal = ""
		tok.Type = tokens.EOF
	case '=':
		if l.PeekChar() == '=' {
			ch := l.ch
			l.ReadChar()
			literal := string(ch) + string(l.ch)
			tok = tokens.Token{Type: tokens.EQ, Literal: literal}

		} else {
			tok = newToken(tokens.ASSIGN, l.ch)
		}
	case '!':
		if l.PeekChar() == '=' {
			ch := l.ch
			l.ReadChar()
			literal := string(ch) + string(l.ch)
			tok = tokens.Token{Type: tokens.NEQ, Literal: literal}

		} else {
			tok = newToken(tokens.BANG, l.ch)
		}
	default:
		if isLetter(l.ch) {
			tok.Literal = l.ReadIdentifier()
			tok.Type = tokens.LookupIdent(tok.Literal)
			return tok
		} else if IsDigit(l.ch) {
			tok.Type = tokens.INT
			tok.Literal = l.ReadNumber()
			return tok

		} else {
			tok = newToken(tokens.ILLEGAL, l.ch)
		}

	}

	l.ReadChar()
	return tok

}

func (l *Lexer) readString() string {
	position := l.Position + 1
	for {
		l.ReadChar()
		if l.ch == '"' || l.ch == 0 {

			break
		}
	}
	return l.Input[position:l.Position]
}

func (l *Lexer) PeekChar() byte {
	if l.Readposition >= len(l.Input) {
		return 0
	} else {
		return l.Input[l.Readposition]
	}

}

func (l *Lexer) ReadIdentifier() string {

	position := l.Position

	for isLetter(l.ch) {

		l.ReadChar()
	}

	return l.Input[position:l.Position]

}

func (l *Lexer) SkipWhiteSpaces() {
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\t' || l.ch == '\r' {
		l.ReadChar()
	}
}

func isLetter(ch byte) bool {

	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'

}

func IsDigit(ch byte) bool {

	return '0' <= ch && ch <= '9'

}

func (l *Lexer) ReadNumber() string {
	position := l.Position
	for IsDigit(l.ch) {
		l.ReadChar()
	}
	return l.Input[position:l.Position]
}

func newToken(token tokens.TokenType, ch byte) tokens.Token {
	return tokens.Token{Type: token, Literal: string(ch)}

}

func New(input string) *Lexer {
	l := &Lexer{Input: input}
	l.ReadChar()
	return l
}

func (l *Lexer) ReadChar() {
	if l.Readposition >= len(l.Input) {
		l.ch = 0
	} else {
		l.ch = l.Input[l.Readposition]
	}
	l.Position = l.Readposition
	l.Readposition += 1
}
