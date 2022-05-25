package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func run(input string, t *testing.T) {
	qry, _ := ParseQuery(input)
	assert.NotEqual(t, qry.TableName, "")
}

func TestParseSelect1(t *testing.T) {
	run("SELECT 1 from xundas", t)
}
func TestParseSelect2(t *testing.T) {
	run("SELECT 1 FROM x WHERE y IN ('a', 'b', 'c') AND o = 2", t)
}
func TestParseCreateTable(t *testing.T) {
	run("CREATE TABLE types (a float(2), b float(49), c NUMERIC(2, 3), d character(4), e char(5), f varchar(6), g character varying(7)) shard set(a, b)", t)
}
