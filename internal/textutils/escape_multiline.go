package textutils

import (
	"bytes"

	"github.com/JohannesKaufmann/html-to-markdown/v2/marker"
)

var newline = []byte{'\n'}
var escape = []byte{'\\'}

func EscapeMultiLine(content []byte) []byte {
	content = bytes.TrimSpace(content)
	content = TrimConsecutiveNewlines(content)
	if len(content) == 0 {
		return content
	}

	parts := marker.SplitFunc(content, func(r rune) bool {
		return r == '\n' || r == marker.MarkerLineBreak
	})

	for i := range parts {
		parts[i] = bytes.TrimSpace(parts[i])
		if len(parts[i]) == 0 {
			parts[i] = escape
		}
	}
	content = bytes.Join(parts, newline)

	return content
}

/*
// TODO: use this optimized function again after integrating the marker.MarkerLineBreak changes

// EscapeMultiLine deals with multiline content inside a link or a heading.
func EscapeMultiLine(content []byte) []byte {
	content = TrimConsecutiveNewlines(content)

	newContent := make([]byte, 0, len(content))

	startNormal := 0
	lineHasContent := false
	for index, char := range content {
		isNewline := char == '\n'
		isSpace := char == ' ' || char == '	'

		isFirstNewline := isNewline && lineHasContent
		isLastNewline := isNewline && !lineHasContent

		if isFirstNewline {
			newContent = append(newContent, content[startNormal:index]...)
			newContent = append(newContent, '\n')

			startNormal = index + 1
			lineHasContent = false

			continue
		} else if isLastNewline {
			newContent = append(newContent, '\\')
			newContent = append(newContent, '\n')

			startNormal = index + 1
			lineHasContent = false
		} else if !isSpace {
			lineHasContent = true
		} else if isSpace && !lineHasContent {
			startNormal = index + 1
		}
	}

	newContent = append(newContent, content[startNormal:]...)

	return newContent
}
*/
