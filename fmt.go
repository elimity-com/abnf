package abnf

import (
	"regexp"
	"strings"
)

func formatRuleName(name string) string {
	var formatted string
	splitHyphen := strings.Split(name, "-")
	for _, part := range splitHyphen {
		formatted += strings.Title(part)
	}
	return formatted
}

func formatFuncComment(comment string) string {
	comment = removeDuplicateSpace(strings.Split(comment, ";")[0])
	splitEqual := strings.SplitAfter(comment, " ")
	return formatRuleName(splitEqual[0]) + strings.Join(splitEqual[1:], "")
}

var space = regexp.MustCompile(`\s+`)

func removeDuplicateSpace(s string) string {
	return space.ReplaceAllString(s, " ")
}

var nl = regexp.MustCompile(`\r?\n`)

func removeNewLines(s string) string {
	return nl.ReplaceAllString(s, "")
}
