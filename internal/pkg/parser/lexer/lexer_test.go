package lexer

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToken_lexNumeric(t *testing.T) {
	tests := []struct {
		number bool
		value  string
	}{
		{
			number: true,
			value:  "105",
		},
		{
			number: true,
			value:  "105 ",
		},
		{
			number: true,
			value:  "123.",
		},
		{
			number: true,
			value:  "123.145",
		},
		{
			number: true,
			value:  "1e5",
		},
		{
			number: true,
			value:  "1.e21",
		},
		{
			number: true,
			value:  "1.1e2",
		},
		{
			number: true,
			value:  "1.1e-2",
		},
		{
			number: true,
			value:  "1.1e+2",
		},
		{
			number: true,
			value:  "1e-1",
		},
		{
			number: true,
			value:  ".1",
		},
		{
			number: true,
			value:  "4.",
		},
		// false tests
		{
			number: false,
			value:  "e4",
		},
		{
			number: false,
			value:  "1..",
		},
		{
			number: false,
			value:  "1ee4",
		},
		{
			number: false,
			value:  " 1",
		},
	}

	for _, test := range tests {
		tok, _, ok := lexNumeric(test.value, cursor{})
		assert.Equal(t, test.number, ok, test.value)
		if ok {
			assert.Equal(t, strings.TrimSpace(test.value), tok.Value, test.value)
		}
	}
}

func TestToken_lexString(t *testing.T) {
	tests := []struct {
		string bool
		value  string
	}{
		{
			string: false,
			value:  "a",
		},
		{
			string: true,
			value:  "'abc'",
		},
		{
			string: true,
			value:  "'a b'",
		},
		{
			string: true,
			value:  "'a' ",
		},
		{
			string: true,
			value:  "'a '' b'",
		},
		// false tests
		{
			string: false,
			value:  "'",
		},
		{
			string: false,
			value:  "",
		},
		{
			string: false,
			value:  " 'foo'",
		},
	}

	for _, test := range tests {
		tok, _, ok := lexString(test.value, cursor{})
		assert.Equal(t, test.string, ok, test.value)
		if ok {
			test.value = strings.TrimSpace(test.value)
			assert.Equal(t, test.value[1:len(test.value)-1], tok.Value, test.value)
		}
	}
}

func TestToken_lexSymbol(t *testing.T) {
	tests := []struct {
		symbol bool
		value  string
	}{
		{
			symbol: true,
			value:  "= ",
		},
		{
			symbol: true,
			value:  "||",
		},
	}

	for _, test := range tests {
		tok, _, ok := lexSymbol(test.value, cursor{})
		assert.Equal(t, test.symbol, ok, test.value)
		if ok {
			test.value = strings.TrimSpace(test.value)
			assert.Equal(t, test.value, tok.Value, test.value)
		}
	}
}

func TestToken_lexIdentifier(t *testing.T) {
	tests := []struct {
		Identifier bool
		input      string
		value      string
	}{
		{
			Identifier: true,
			input:      "a",
			value:      "a",
		},
		{
			Identifier: true,
			input:      "abc",
			value:      "abc",
		},
		{
			Identifier: true,
			input:      "abc ",
			value:      "abc",
		},
		{
			Identifier: true,
			input:      `" abc "`,
			value:      ` abc `,
		},
		{
			Identifier: true,
			input:      "a9$",
			value:      "a9$",
		},
		{
			Identifier: true,
			input:      "userName",
			value:      "username",
		},
		{
			Identifier: true,
			input:      `"userName"`,
			value:      "userName",
		},
		// false tests
		{
			Identifier: false,
			input:      `"`,
		},
		{
			Identifier: false,
			input:      "_sadsfa",
		},
		{
			Identifier: false,
			input:      "9sadsfa",
		},
		{
			Identifier: false,
			input:      " abc",
		},
	}

	for _, test := range tests {
		tok, _, ok := lexIdentifier(test.input, cursor{})
		assert.Equal(t, test.Identifier, ok, test.input)
		if ok {
			assert.Equal(t, test.value, tok.Value, test.input)
		}
	}
}

func TestToken_lexKeyword(t *testing.T) {
	tests := []struct {
		keyword bool
		value   string
	}{
		{
			keyword: true,
			value:   "select ",
		},
		{
			keyword: true,
			value:   "from",
		},
		{
			keyword: true,
			value:   "as",
		},
		{
			keyword: true,
			value:   "SELECT",
		},
		{
			keyword: true,
			value:   "into",
		},
		// false tests
		{
			keyword: false,
			value:   " into",
		},
		{
			keyword: false,
			value:   "flubbrety",
		},
	}

	for _, test := range tests {
		tok, _, ok := lexKeyword(test.value, cursor{})
		assert.Equal(t, test.keyword, ok, test.value)
		if ok {
			test.value = strings.TrimSpace(test.value)
			assert.Equal(t, strings.ToLower(test.value), tok.Value, test.value)
		}
	}
}

func TestLex(t *testing.T) {
	tests := []struct {
		input  string
		Tokens []Token
		err    error
	}{
		{
			input: "select a",
			Tokens: []Token{
				{
					Loc:   Location{Col: 0, Line: 0},
					Value: string(keywordSelect),
					Kind:  KindKeyword,
				},
				{
					Loc:   Location{Col: 7, Line: 0},
					Value: "a",
					Kind:  KindIdentifier,
				},
			},
		},
		{
			input: "select true",
			Tokens: []Token{
				{
					Loc:   Location{Col: 0, Line: 0},
					Value: string(keywordSelect),
					Kind:  KindKeyword,
				},
				{
					Loc:   Location{Col: 7, Line: 0},
					Value: "true",
					Kind:  KindBool,
				},
			},
		},
		{
			input: "select 1",
			Tokens: []Token{
				{
					Loc:   Location{Col: 0, Line: 0},
					Value: string(keywordSelect),
					Kind:  KindKeyword,
				},
				{
					Loc:   Location{Col: 7, Line: 0},
					Value: "1",
					Kind:  KindNumeric,
				},
			},
			err: nil,
		},
		{
			input: "select 'foo' || 'bar';",
			Tokens: []Token{
				{
					Loc:   Location{Col: 0, Line: 0},
					Value: string(keywordSelect),
					Kind:  KindKeyword,
				},
				{
					Loc:   Location{Col: 7, Line: 0},
					Value: "foo",
					Kind:  KindString,
				},
				{
					Loc:   Location{Col: 13, Line: 0},
					Value: string(SymbolConcat),
					Kind:  KindSymbol,
				},
				{
					Loc:   Location{Col: 16, Line: 0},
					Value: "bar",
					Kind:  KindString,
				},
				{
					Loc:   Location{Col: 21, Line: 0},
					Value: string(SymbolSemicolon),
					Kind:  KindSymbol,
				},
			},
			err: nil,
		},
		{
			input: "CREATE TABLE u (id INT, name TEXT)",
			Tokens: []Token{
				{
					Loc:   Location{Col: 0, Line: 0},
					Value: string(keywordCreate),
					Kind:  KindKeyword,
				},
				{
					Loc:   Location{Col: 7, Line: 0},
					Value: string(keywordTable),
					Kind:  KindKeyword,
				},
				{
					Loc:   Location{Col: 13, Line: 0},
					Value: "u",
					Kind:  KindIdentifier,
				},
				{
					Loc:   Location{Col: 15, Line: 0},
					Value: "(",
					Kind:  KindSymbol,
				},
				{
					Loc:   Location{Col: 16, Line: 0},
					Value: "id",
					Kind:  KindIdentifier,
				},
				{
					Loc:   Location{Col: 19, Line: 0},
					Value: "int",
					Kind:  KindKeyword,
				},
				{
					Loc:   Location{Col: 22, Line: 0},
					Value: ",",
					Kind:  KindSymbol,
				},
				{
					Loc:   Location{Col: 24, Line: 0},
					Value: "name",
					Kind:  KindIdentifier,
				},
				{
					Loc:   Location{Col: 29, Line: 0},
					Value: "text",
					Kind:  KindKeyword,
				},
				{
					Loc:   Location{Col: 33, Line: 0},
					Value: ")",
					Kind:  KindSymbol,
				},
			},
		},
		{
			input: "insert into users Values (105, 233)",
			Tokens: []Token{
				{
					Loc:   Location{Col: 0, Line: 0},
					Value: string(keywordInsert),
					Kind:  KindKeyword,
				},
				{
					Loc:   Location{Col: 7, Line: 0},
					Value: string(keywordInto),
					Kind:  KindKeyword,
				},
				{
					Loc:   Location{Col: 12, Line: 0},
					Value: "users",
					Kind:  KindIdentifier,
				},
				{
					Loc:   Location{Col: 18, Line: 0},
					Value: string(keywordValues),
					Kind:  KindKeyword,
				},
				{
					Loc:   Location{Col: 25, Line: 0},
					Value: "(",
					Kind:  KindSymbol,
				},
				{
					Loc:   Location{Col: 26, Line: 0},
					Value: "105",
					Kind:  KindNumeric,
				},
				{
					Loc:   Location{Col: 30, Line: 0},
					Value: ",",
					Kind:  KindSymbol,
				},
				{
					Loc:   Location{Col: 32, Line: 0},
					Value: "233",
					Kind:  KindNumeric,
				},
				{
					Loc:   Location{Col: 36, Line: 0},
					Value: ")",
					Kind:  KindSymbol,
				},
			},
			err: nil,
		},
		{
			input: "SELECT id FROM users;",
			Tokens: []Token{
				{
					Loc:   Location{Col: 0, Line: 0},
					Value: string(keywordSelect),
					Kind:  KindKeyword,
				},
				{
					Loc:   Location{Col: 7, Line: 0},
					Value: "id",
					Kind:  KindIdentifier,
				},
				{
					Loc:   Location{Col: 10, Line: 0},
					Value: string(keywordFrom),
					Kind:  KindKeyword,
				},
				{
					Loc:   Location{Col: 15, Line: 0},
					Value: "users",
					Kind:  KindIdentifier,
				},
				{
					Loc:   Location{Col: 20, Line: 0},
					Value: ";",
					Kind:  KindSymbol,
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		tokens, err := Lex(test.input)
		assert.Equal(t, test.err, err, test.input)

		msg := fmt.Sprintf("%s: %s", test.input, err)

		assert.Equal(t, len(test.Tokens), len(tokens), msg)

		for i, tok := range tokens {
			assert.Equal(t, &test.Tokens[i], tok, test.input)
		}
	}
}

func BenchmarkLex(b *testing.B) {
	sqls := map[string]string{
		"create table": "CREATE TABLE u (id INT, name TEXT)",
		"s a":          "select a",
		"s true":       "select true",
		"s 1":          "select 1",
		"s 1.1":        "select 1.1",
		"s or":         "select 'foo' || 'bar';",
		"ins":          "insert into users Values (105, 233)",
		"s id ":        "SELECT id FROM users;",
		"s big": `select
		"acccount_manager"."provider".*,
		array_agg(
			"acccount_manager"."provider_domain"."domain_name"
		) as "domains",
		array_agg(
			"acccount_manager"."provider_origin"."origin"
		) as "origins"
	from
		"acccount_manager"."login_rollout"
		join "acccount_manager"."provider_login" on "acccount_manager"."login_rollout"."login_id" = "acccount_manager"."provider_login"."login_id"
		join "acccount_manager"."provider" on "acccount_manager"."provider_login"."provider_id" = "acccount_manager"."provider"."provider_id"
		left outer join "acccount_manager"."provider_domain" on (
			"acccount_manager"."provider"."provider_id" = "acccount_manager"."provider_domain"."provider_id"
			and "acccount_manager"."provider_domain"."enabled"
		)
		left outer join "acccount_manager"."provider_origin" on (
			"acccount_manager"."provider"."provider_id" = "acccount_manager"."provider_origin"."provider_id"
			and "acccount_manager"."provider_origin"."enabled"
		)
	where
		(
			lower(
				"acccount_manager"."login_rollout"."login"
			) = lower($1)
			and "acccount_manager"."login_rollout"."enabled"
		)
	group by
		"acccount_manager"."provider"."provider_id",
		"acccount_manager"."provider"."code",
		"acccount_manager"."provider"."priority",
		"acccount_manager"."provider"."enabled",
		"acccount_manager"."provider"."is_corporate",
		"acccount_manager"."provider"."created_at",
		"acccount_manager"."provider"."updated_at"
	order by
		"acccount_manager"."provider"."priority" desc`,
	}

	for k, sql := range sqls {
		// if k == "s bolado" {
		// 	continue
		// }
		b.Run(k, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = Lex(sql)
			}
		})
	}
}
