package parser

import "strings"

func getConditionsWithWhereClause(query string) []Condition {
	strConditions := strings.Split(after(query, "where"), "and")
	var conditions []Condition

	for _, v := range strConditions {
		if strings.Contains(v, "=") {
			val := strings.Split(v, "=")
			cond := Condition{
				Field: val[0],
				Value: val[1],
			}
			conditions = append(conditions, cond)
		}
	}

	return conditions
}

func getConditionsWithInsertClause(query string) []Condition {
	fields := strings.Split(between(query, "(", ")"), ",")
	strValues := after(query, "values")
	values := strings.Split(strings.ReplaceAll(strings.ReplaceAll(strValues, "(", ""), ")", ""), ",")
	var conditions []Condition

	for i, v := range values {
		cond := Condition{
			Field: fields[i],
			Value: v,
		}
		conditions = append(conditions, cond)
	}

	return conditions
}
