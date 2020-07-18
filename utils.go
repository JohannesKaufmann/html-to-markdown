package md

import (
	"regexp"
	"strings"
	"unicode"
)

func TrimpLeadingSpaces(text string) string {
	parts := strings.Split(text, "\n")
	for i := range parts {
		b := []byte(parts[i])

		var spaces int
		for i := 0; i < len(b); i++ {
			if unicode.IsSpace(rune(b[i])) {
				if b[i] == '	' {
					spaces = spaces + 4
				} else {
					spaces++
				}
				continue
			}

			// this seems to be a list item
			if b[i] == '-' {
				break
			}

			// this seems to be a code block
			if spaces >= 4 {
				break
			}

			// remove the space characters from the string
			b = b[i:]
			break
		}
		parts[i] = string(b)
	}

	return strings.Join(parts, "\n")
}

func TrimTrailingSpaces(text string) string {
	parts := strings.Split(text, "\n")
	for i := range parts {
		parts[i] = strings.TrimRightFunc(parts[i], func(r rune) bool {
			return unicode.IsSpace(r)
		})

	}

	return strings.Join(parts, "\n")
}

// The same as `multipleNewLinesRegex`, but applies to escaped new lines inside a link `\n\`
var multipleNewLinesInLinkRegex = regexp.MustCompile(`(\n\\){1,}`) // `([\n\r\s]\\)`

func EscapeMultiLine(content string) string {
	content = strings.TrimSpace(content)
	content = strings.Replace(content, "\n", `\`+"\n", -1)

	content = multipleNewLinesInLinkRegex.ReplaceAllString(content, "\n\\\n\\")

	return content
}
