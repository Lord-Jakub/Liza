package main

import (
	"fmt"

	"lizalang/lexer"
	//"lizalang/token"
)

func main() {
	code := `func main(){
		string helloworld = "Hello World"
		print (helloworld)
	}`
	code = string(append([]byte(code), 0))
	lex := lexer.New(code, "nil")
	lex.NextToken()
	for _, tok := range lex.Tokens {
		fmt.Printf("%s:%s\n", tok.Type, tok.Value)
	}
}
