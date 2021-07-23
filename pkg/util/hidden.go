package util

import (
	"regexp"
)

func HiddePass(url string) string {
	var re = regexp.MustCompile(`(?m)(\/\/\w+\:)(.*)(@)`)
	var substitution = "${1}******${3}"

	return re.ReplaceAllString(url, substitution)
}
