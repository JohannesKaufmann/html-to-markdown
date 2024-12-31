package textutils

// EscapeMultiLine deals with multiline content inside a link or a heading.
func EscapeMultiLine(content []byte) []byte {
	content = Alternative_TrimConsecutiveNewlines(content)

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
