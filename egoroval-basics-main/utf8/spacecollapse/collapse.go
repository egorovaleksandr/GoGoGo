//go:build !solution

package spacecollapse

import (
	"unicode"
)

func CollapseSpaces(input string) string {
	output := []rune{}
	inSpaces := false
	for _, r := range input {
		if unicode.IsSpace(r) {
			if inSpaces {
				continue
			}
			inSpaces = true
			output = append(output, ' ')
		} else {
			inSpaces = false
			output = append(output, r)
		}
	}
	return string(output)
}
