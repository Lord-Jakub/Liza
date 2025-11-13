package lexer

import (
	"fmt"
	"slices"
	"strconv"

	"lizalang/token"
	"lizalang/utils"
)

type Lexer struct {
	Code      []byte
	Pos       int
	Line      int
	CurChar   byte
	CharAfter byte
	Tokens    []token.Token
	File      string
	Errors    []error
}

func (lexer *Lexer) NextChar() {
	lexer.Pos++
	lexer.CurChar = lexer.CharAfter
	if lexer.Pos+1 < len(lexer.Code) {
		lexer.CharAfter = lexer.Code[lexer.Pos+1]
	}
}

func New(code string, file string) *Lexer {
	return &Lexer{
		[]byte(code),
		0,
		1,
		[]byte(code)[0],
		[]byte(code)[1],
		make([]token.Token, 0),
		file,
		make([]error, 0),
	}
}

func (lexer *Lexer) NextToken() {
	switch {
	case lexer.CurChar == 0:
		lexer.NewToken(token.EOF, 0)
		break
	case utils.IsLetter(lexer.CurChar):
		lexer.handleIdOrKeyword()
		break
	case utils.IsDigit(lexer.CurChar):
		err := lexer.handleNumber()
		if err != nil {
			lexer.Errors = append(lexer.Errors, err)
		}
		break
	case lexer.CurChar == '<' && lexer.CharAfter == '=':
		lexer.NewToken(token.LessThanOrEqual, "<=")
		lexer.NextChar()
		break
	case lexer.CurChar == '>' && lexer.CharAfter == '=':
		lexer.NewToken(token.MoreThanOrEqual, ">=")
		lexer.NextChar()
		break
	case lexer.CurChar == '!' && lexer.CharAfter == '=':
		lexer.NewToken(token.NotEqual, "!=")
		lexer.NextChar()
		break
	case lexer.CurChar == '=' && lexer.CharAfter == '=':
		lexer.NewToken(token.DoubleEqual, "==")
		lexer.NextChar()
		break
	case lexer.CurChar == '"':
		lexer.handleString()
		break
	case lexer.CurChar == '\n':
		lexer.NewToken(token.NewInstruction, lexer.CurChar)
		lexer.Line++
		if lexer.CharAfter == '\t' {
			lexer.NextChar()
		}
		break
	case lexer.CurChar == ' ' || lexer.CurChar == '\t':
		break
	default:
		if oneCharToken, ok := token.SymbolMap[lexer.CurChar]; ok {
			lexer.NewToken(oneCharToken, lexer.CurChar)
			break
		}
		lexer.NewToken(token.Invalid, lexer.CurChar)
		lexer.Errors = append(lexer.Errors, fmt.Errorf("Invalid character %s on a line %d", lexer.CurChar, lexer.Line))
		break
	}
	if lexer.Tokens[len(lexer.Tokens)-1].Type != token.EOF {
		lexer.NextChar()
		lexer.NextToken()
	}
}

func (lexer *Lexer) NewToken(tokentype token.TokenType, value any) {
	tok := token.Token{
		tokentype,
		value,
		lexer.Line,
		lexer.File,
	}
	lexer.Tokens = append(lexer.Tokens, tok)
}

func (lexer *Lexer) handleIdOrKeyword() {
	str := ""
	for utils.IsLetter(lexer.CurChar) {
		str = string(append([]byte(str), lexer.CurChar))
		if utils.IsLetter(lexer.CharAfter) {
			lexer.NextChar()
		} else {
			break
		}
	}
	if slices.Contains(token.KeyWords, str) {
		lexer.NewToken(token.Keyword, str)
	} else {
		lexer.NewToken(token.Identifier, str)
	}
}

func (lexer *Lexer) handleNumber() error {
	num := ""
	hasDot := false
	for utils.IsLetter(lexer.CurChar) || lexer.CurChar == '.' {
		num = string(append([]byte(num), lexer.CurChar))
		if lexer.CurChar == '.' {
			hasDot = true
		}
		if utils.IsDigit(lexer.CharAfter) || lexer.CharAfter == '.' {
			lexer.NextChar()
		} else {
			break
		}
	}
	if hasDot {
		if floatNum, err := strconv.ParseFloat(num, 64); err == nil {
			lexer.NewToken(token.Float, floatNum)
		} else {
			lexer.NewToken(token.Invalid, "NaN")
			return fmt.Errorf("Error at line %d: %s is not a number", lexer.Line, num)
		}
	} else {
		intNum, err := strconv.ParseInt(num, 10, 64)
		if err != nil {
			return fmt.Errorf("Error at line %d: %s is not a number", lexer.Line, num)
		}
		lexer.NewToken(token.Int, intNum)
	}
	return nil
}

func (lexer *Lexer) handleString() {
	lexer.NextChar()
	str := ""
	for lexer.CurChar != '"' {
		str = string(append([]byte(str), lexer.CurChar))
		lexer.NextChar()
	}
	lexer.NewToken(token.String, str)
}
