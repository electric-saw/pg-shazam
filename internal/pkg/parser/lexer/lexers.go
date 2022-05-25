package lexer

import "strings"

func lexNumeric(source string, srcCur cursor) (*Token, cursor, bool) {

	periodFound := false
	expFound := false
	cur := srcCur

	for ; cur.pointer < uint(len(source)); cur.pointer++ {
		c := source[cur.pointer]
		cur.loc.Col++

		isDigit := c >= '0' && c <= '9'
		isPeriod := c == '.'
		isExp := c == 'e'

		if cur.pointer == srcCur.pointer {
			if !isDigit && !isPeriod {
				return nil, srcCur, false
			}

			periodFound = isPeriod
			continue
		}

		if isPeriod {
			if periodFound {
				return nil, srcCur, false
			}

			periodFound = true
			continue
		}

		if isExp {
			if expFound {
				return nil, srcCur, false
			}

			periodFound = true
			expFound = true

			if cur.pointer == uint(len(source)-1) {
				return nil, srcCur, false
			}

			cNext := source[cur.pointer+1]

			if cNext == '+' || cNext == '-' {
				cur.pointer++
				cur.loc.Col++
			}

			continue

		}

		if !isDigit {
			break
		}
	}

	if cur.pointer == srcCur.pointer {
		return nil, srcCur, false
	}

	return &Token{
		Kind:  KindNumeric,
		Value: source[srcCur.pointer:cur.pointer],
		Loc:   srcCur.loc,
	}, cur, true
}

func lexCharDelimited(source string, srcCur cursor, delimiter byte) (*Token, cursor, bool) {
	cur := srcCur

	if len(source) == 0 {
		return nil, srcCur, false
	}

	if source[cur.pointer] != delimiter {
		return nil, srcCur, false
	}

	cur.pointer++
	cur.loc.Col++

	for ; cur.pointer < uint(len(source)); cur.pointer++ {
		c := source[cur.pointer]

		if c == delimiter {
			if cur.pointer+1 >= uint(len(source)) || source[cur.pointer+1] != delimiter {
				cur.pointer++
				cur.loc.Col++

				return &Token{
					Kind:  KindString,
					Value: source[srcCur.pointer+1 : cur.pointer-1],
					Loc:   srcCur.loc,
				}, cur, true
			}

			cur.pointer++
			cur.loc.Col++

			if cur.pointer == '\n' {
				cur.loc.Line++
				cur.loc.Col = 0
			}

		}
		cur.loc.Col++
	}

	return nil, srcCur, false
}

func lexString(source string, srcCur cursor) (*Token, cursor, bool) {
	return lexCharDelimited(source, srcCur, '\'')
}

// lex sql comment -- or /*
func lexComment(source string, srcCur cursor) (*Token, cursor, bool) {
	cur := srcCur

	nextChar := '-'
	if source[cur.pointer] == '/' {
		nextChar = '*'
	}

	if source[cur.pointer] == '-' || source[cur.pointer] == '/' {
		cur.pointer++
		cur.loc.Col++

		if cur.pointer == uint(len(source)) {
			return nil, srcCur, false
		}

		if source[cur.pointer] == byte(nextChar) {
			cur.pointer++
			cur.loc.Col++

			for ; cur.pointer < uint(len(source)); cur.pointer++ {
				c := source[cur.pointer]
				cur.loc.Col++

				if c == '\n' {
					return &Token{
						Kind:  KindComment,
						Value: source[srcCur.pointer:cur.pointer],
						Loc:   srcCur.loc,
					}, cur, true
				}
			}

			return &Token{
				Kind:  KindComment,
				Value: source[srcCur.pointer:cur.pointer],
				Loc:   srcCur.loc,
			}, cur, true
		}

		return nil, srcCur, false
	}

	return nil, srcCur, false
}

func lexWhitespace(source string, srcCur cursor) (*Token, cursor, bool) {
	cur := srcCur

	for ; cur.pointer < uint(len(source)); cur.pointer++ {
		c := source[cur.pointer]

		if c != ' ' && c != '\t' && c != '\n' && c != '\r' {
			break
		}

		if c == '\n' {
			cur.loc.Line++
			cur.loc.Col = 0
		} else {
			cur.loc.Col++
		}
	}

	if cur.pointer == srcCur.pointer {
		return nil, srcCur, false
	}

	return &Token{
		Kind:  KindWhitespace,
		Value: " ",
		Loc:   srcCur.loc,
	}, cur, true
}

func lexIdentifier(source string, ic cursor) (*Token, cursor, bool) {
	if token, newCursor, ok := lexCharDelimited(source, ic, '"'); ok {
		return token, newCursor, true
	}

	cur := ic

	c := source[cur.pointer]
	// Other characters count too, big ignoring non-ascii for now
	isAlphabetical := (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
	if !isAlphabetical {
		return nil, ic, false
	}
	cur.pointer++
	cur.loc.Col++

	for ; cur.pointer < uint(len(source)); cur.pointer++ {
		c = source[cur.pointer]

		// Other characters count too, big ignoring non-ascii for now
		isAlphabetical := (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
		isNumeric := c >= '0' && c <= '9'
		if isAlphabetical || isNumeric || c == '$' || c == '_' {
			cur.loc.Col++
			continue
		}

		break
	}

	value := source[ic.pointer:cur.pointer]

	if len(value) == 0 {
		return nil, ic, false
	}

	return &Token{
		Value: strings.ToLower(value),
		Loc:   ic.loc,
		Kind:  KindIdentifier,
	}, cur, true
}
