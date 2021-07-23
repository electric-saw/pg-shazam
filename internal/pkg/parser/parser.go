package parser

import (
	"strings"
)

type Operation int8

const (
	Select Operation = iota
	Insert
	Update
	Delete
	CreateTable
	Set
)

type Condition struct {
	Field string
	Value string
}

type Query struct {
	Operation    Operation
	TableName    string
	Conditions   []Condition
	Shards       []string
	DDLOperation bool
	QueryString  string
}

func ParseQuery(query string) Query {
	norQuery := normalize(query)
	op, isDDLOperation := howIsOperation(norQuery)
	// TODO: adjust error handling
	table, _ := howIsMainTable(norQuery, op)

	if isDDLOperation {
		return Query{
			Operation:    op,
			TableName:    table,
			Shards:       getColumnsShard(norQuery),
			Conditions:   howIsConditionsEquals(norQuery, op),
			DDLOperation: true,
			QueryString:  removeShardInQuery(query),
		}
	} else {
		return Query{
			Operation:    op,
			TableName:    table,
			Conditions:   howIsConditionsEquals(norQuery, op),
			DDLOperation: false,
			QueryString:  query,
		}
	}
}

func howIsOperation(query string) (Operation, bool) {
	var operation Operation
	isDDLOperation := false
	query = strings.ToLower(query)

	if strings.HasPrefix(query, "select") {
		operation = Operation(Select)
	} else if strings.HasPrefix(query, "update") {
		operation = Operation(Update)
	} else if strings.HasPrefix(query, "insert into") {
		operation = Operation(Insert)
	} else if strings.HasPrefix(query, "delete from") {
		operation = Operation(Delete)
	} else if strings.HasPrefix(query, "create table") {
		operation = Operation(CreateTable)
		isDDLOperation = true
	} else if strings.HasPrefix(query, "create") {
		isDDLOperation = true
	} else if strings.HasPrefix(query, "drop") {
		isDDLOperation = true
	} else if strings.HasPrefix(query, "alter") {
		isDDLOperation = true
	} else if strings.HasPrefix(query, "truncate") {
		isDDLOperation = true
	} else if strings.HasPrefix(query, "set") {
		operation = Operation(Set)
	}
	return operation, isDDLOperation
}

func howIsMainTable(query string, op Operation) (string, error) {
	if op == Operation(Select) {
		return getWordAfter(query, "from")
	} else if op == Operation(Update) {
		return getWordAfter(query, "update")
	} else if op == Operation(Insert) {
		return getWordAfter(query, "into")
	} else if op == Operation(Delete) {
		return getWordAfter(query, "from")
	} else if op == Operation(CreateTable) {
		return getWordAfter(query, "table")
	}

	return "", nil
}

func howIsConditionsEquals(query string, op Operation) []Condition {
	if op == Operation(Select) {
		return getConditionsWithWhereClause(query)
	} else if op == Operation(Update) {
		return getConditionsWithWhereClause(query)
	} else if op == Operation(Insert) {
		return getConditionsWithInsertClause(query)
	} else if op == Operation(Delete) {
		return getConditionsWithWhereClause(query)
	}
	return nil
}
