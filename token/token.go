package token

type Token struct {
	Type  TokenType
	Value any
	Line  int
	File  string
}
type TokenType string

const (
	Invalid         = "INVALID"
	EOF             = "EOF"
	OpenParen       = "OPENPAREN"
	CloseParen      = "CLOSEPAREN"
	OpenBrace       = "OPENBRACE"
	CloseBrace      = "CLOSEBRACE"
	Plus            = "PLUS"
	Minus           = "MINUS"
	Multiply        = "MULTIPLY"
	Divide          = "DIVIDE"
	Backslash       = "BACKSLASH"
	NewInstruction  = "NEWINSTRUCTION"
	SingleQuote     = "SINGLEQUOTE"
	DoubleQuote     = "DOUBLEQUOTE"
	Equal           = "EQUAL"
	Not             = "NOT"
	Comma           = "COMMA"
	LessThan        = "LESSTHAN"
	LessThanOrEqual = "LESSTHANOREQUAL"
	MoreThan        = "MORETHAN"
	MoreThanOrEqual = "MORETHANOREQUAL"
	Int             = "INT"
	String          = "STRING"
	Identifier      = "IDENTIFIER"
	Keyword         = "KEYWORD"
	DoubleEqual     = "DOUBLEEQUAL"
	NotEqual        = "NOTEQUAL"
	Float           = "FLOAT"
	OpenBracket     = "OPENBRACKET"
	CloseBracket    = "CLOSEBRACKET"
)

var SymbolMap = map[rune]TokenType{
	'(':  OpenParen,
	')':  CloseParen,
	'{':  OpenBrace,
	'}':  CloseBrace,
	'+':  Plus,
	'-':  Minus,
	'*':  Multiply,
	'/':  Divide,
	'\\': Backslash,
	'=':  Equal,
	'!':  Not,
	',':  Comma,
	'<':  LessThan,
	'>':  MoreThan,
	';':  NewInstruction,
	'[':  OpenBracket,
	']':  CloseBracket,
}
var KeyWords = []string{"if", "else", "for", "func", "return", "string", "int", "float", "bool", "void"}
