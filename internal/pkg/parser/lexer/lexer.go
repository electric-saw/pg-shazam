//go:generate go run ./gen/
package lexer

import (
	"fmt"
)

type Location struct {
	Line int
	Col  int
}

var lexers = []lexer{
	lexWhitespace,
	lexNumeric,
	lexString,
	lexBool,
	lexComment,
	lexKeyword,
	lexSymbol,
	lexIdentifier,
}

type lexer func(string, cursor) (*Token, cursor, bool)

func Lex(source string) ([]*Token, error) {
	tokens := []*Token{}
	cur := cursor{}

lex:
	for cur.pointer < uint(len(source)) {
		for _, lexer := range lexers {
			if token, newCur, ok := lexer(source, cur); ok {
				cur = newCur
				if token != nil {
					// fmt.Println(source)
					// fmt.Print(strings.Repeat(" ", int(token.Loc.Col)))
					// fmt.Print("^")
					// miolo := len(token.Value) - 2
					// if token.Kind == KindString {
					// 	miolo += 2
					// }

					// if miolo > 0 {
					// 	fmt.Print(strings.Repeat("-", miolo))
					// }
					// if len(token.Value) > 1 {
					// 	fmt.Print("^")
					// }
					// fmt.Printf(" %d:%d - %s\n", token.Loc.Col, cur.loc.Col, token.Value)

					if token.Kind == KindWhitespace {
						continue lex
					}

					tokens = append(tokens, token)

				}

				continue lex
			}
		}
		hint := ""
		if len(tokens) > 0 {
			hint = fmt.Sprintf("after %s", tokens[len(tokens)-1].Value)
		}
		return nil, fmt.Errorf("unexpected character '%s' at %d:%d %s", string(source[cur.pointer]), cur.loc.Line, cur.loc.Col, hint)

	}
	return tokens, nil
}
