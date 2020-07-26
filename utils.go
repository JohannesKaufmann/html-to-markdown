package md

import (
	"bytes"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func CollectText(n *html.Node) string {
	text := &bytes.Buffer{}
	collectText(n, text)
	return text.String()
}
func collectText(n *html.Node, buf *bytes.Buffer) {

	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collectText(c, buf)
	}
}

// always have a space to the side to recognize the delimiter
func AddSpaceIfNessesary(selec *goquery.Selection, text string) string {

	var prev string
	if len(selec.Nodes) > 0 {
		node := selec.Nodes[0]

		if node.PrevSibling != nil {
			for node = node.PrevSibling; node != nil; node = node.PrevSibling {
				prev = CollectText(node)

				// if the content is empty, try our luck with the next node
				if strings.TrimSpace(prev) != "" {
					break
				}
			}
		}
	}

	lastChar, size := utf8.DecodeLastRuneInString(prev)
	if size > 0 && !unicode.IsSpace(lastChar) {
		text = " " + text
	}

	// - - - - - - - - - - - - - - - - - - - //

	var next string
	if len(selec.Nodes) > 0 {
		node := selec.Nodes[len(selec.Nodes)-1]

		if node.NextSibling != nil {
			for node = node.NextSibling; node != nil; node = node.NextSibling {
				next = CollectText(node)

				// if the content is empty, try our luck with the next node
				if strings.TrimSpace(next) != "" {

					// Right now, this function AddSpaceIfNessesary is used for `a`,
					// `strong`, `b`, `i` and `em`.
					// Don't add another space if the other element is going to add a
					// space already.
					s := &goquery.Selection{Nodes: []*html.Node{node}}
					name := goquery.NodeName(s)

					if name == "a" || name == "strong" || name == "b" || name == "i" || name == "em" {
						next = " "
					}

					break
				}
			}

		}
	}

	firstChar, size := utf8.DecodeRuneInString(next)
	if size > 0 && !unicode.IsSpace(firstChar) && !unicode.IsPunct(firstChar) {
		text += " "
	}

	return text
}

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
