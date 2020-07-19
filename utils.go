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

// Cal can be passed the content of a code block and it returns
// how many fence characters (` or ~) should be used.
//
// This is useful if the html content includes the same fence characters
// for example ```
// -> https://stackoverflow.com/a/49268657
func CalculateCodeFence(fenceChar rune, content string) string {
	var occurrences []int

	var charsTogether int
	for _, char := range content {
		// we encountered a fence character, now count how many
		// are directly afterwards
		if char == fenceChar {
			charsTogether++
		} else if charsTogether != 0 {
			occurrences = append(occurrences, charsTogether)
			charsTogether = 0
		}
	}

	// if the last element in the content was a fenceChar
	if charsTogether != 0 {
		occurrences = append(occurrences, charsTogether)
	}

	repeat := findMax(occurrences)

	// the outer fence block always has to have
	// at least one character more than any content inside
	repeat++

	// you have to have at least three fence characters
	// to be recognized as a code block
	if repeat < 3 {
		repeat = 3
	}

	return strings.Repeat(string(fenceChar), repeat)
}

func findMax(a []int) (max int) {
	for i, value := range a {
		if i == 0 {
			max = a[i]
		}

		if value > max {
			max = value
		}
	}
	return max
}
