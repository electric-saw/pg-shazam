package parser

import (
	"regexp"
	"strings"
)

func getColumnsShard(query string) []string {
	columns := between(query, "shard set(", ")")
	var result []string
	for _, column := range strings.Split(columns, ",") {
		result = append(result, strings.TrimSpace(column))
	}
	return result
}

func removeShardInQuery(query string) string {
	r := regexp.MustCompile(`(?im)[,]{0,1}\s+shard set\([\n\s\w\,]+\)`)
	return r.ReplaceAllString(query, "")

	// tableLines := strings.Split(query, ",")
	// var newQuery []string

	// for _, v := range tableLines {
	// 	if !strings.Contains(v, "shard set") {
	// 		newQuery = append(newQuery, v)
	// 	}
	// }

	// return strings.Join(newQuery, ",")
}
