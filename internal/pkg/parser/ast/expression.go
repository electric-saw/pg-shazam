package ast

import "github.com/electric-saw/pg-shazam/internal/pkg/parser/lexer"

type expressionKind uint

const (
	literalKind expressionKind = iota
)

type expression struct {
	literal *lexer.Token
	Kind    expressionKind
}
