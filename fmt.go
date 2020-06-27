package abnf

import "strings"

func formatRuleName(name string) string {
	var formatted string
	splitHyphen := strings.Split(name, "-")
	for _, part := range splitHyphen {
		formatted += strings.Title(part)
	}
	return formatted
}
