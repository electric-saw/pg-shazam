package parser

type ParseError struct {
	msg string
}

func NewParseError(err string) *ParseError {
	return &ParseError{
		msg: err,
	}
}

func (p *ParseError) Error() string {
	return p.msg
}
