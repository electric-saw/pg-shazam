package parser

import (
	"strings"
)

func after(value string, a string) string {
	// Get substring after a string.
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:]
}

func normalize(txt string) string {
	return strings.ReplaceAll(strings.ToLower(strings.Trim(txt, " ")), "\n", " ")
}

func between(value string, a string, b string) string {
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value[posFirst:], b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	posLastAdjusted := posFirst + posLast

	return value[posFirstAdjusted:posLastAdjusted]
}

func getWordAfter(query string, word string) (string, error) {
	words := strings.Split(query, " ")

	for i, v := range words {
		if v == word {
			return words[i+1], nil
		}
	}

	return "", NewParseError("String not found")
}
