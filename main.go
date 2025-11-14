package main

import (
	"fmt"

	"lizalang/lexer"
	//"lizalang/token"
)

func main() {
	code := `func main(){
		string helloworld = "\"Hello\tworld\"\n
	after two new lines"
		print (helloworld)
	int i = 0
	.
	"
	}`
	code = string(append([]byte(code), 0))
	lex := lexer.New(code, "nil")
	lex.Lex()
	for _, tok := range lex.Tokens {
		if _, ok := tok.Value.(rune); ok {
			fmt.Printf("%s:%s\n", tok.Type, string(tok.Value.(rune)))
		} else {
			fmt.Printf("%s:%s\n", tok.Type, tok.Value)
		}
	}
	for _, err := range lex.Errors {
		fmt.Println(err)
	}
}
