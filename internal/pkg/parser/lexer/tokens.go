package lexer

type TokenKind int

const (
	KindKeyword TokenKind = iota
	KindSymbol
	KindString
	KindNumeric
	KindBool
	KindNull
	KindComment
	KindIdentifier
	KindWhitespace
	KindEOF
)

type Token struct {
	Value string
	Kind  TokenKind
	Loc   Location
}

type cursor struct {
	pointer uint
	loc     Location
}

func (t *Token) Equals(other *Token) bool {
	return t.Value == other.Value && t.Kind == other.Kind
}
