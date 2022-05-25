package ast

import "github.com/electric-saw/pg-shazam/internal/pkg/parser/lexer"

type AstKind uint

const (
	KindUnknown AstKind = iota
	AlterKind
	CreateKind
	DeleteKind
	InsertKind
	SelectKind
	UpdateKind
)

type Statement struct {
	Kind AstKind
	// Alter *AlterStatement
	// Create *CreateStatement
	// Delete *DeleteStatement
	// Insert *InsertStatement
	// Select *SelectStatement
	// Update *UpdateStatement

}

type InsertStatement struct {
	table  lexer.Token
	values *[]*expression
}
