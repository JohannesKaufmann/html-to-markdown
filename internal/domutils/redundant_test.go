package domutils

import (
	"testing"

	"github.com/JohannesKaufmann/dom"
	"github.com/JohannesKaufmann/html-to-markdown/v2/internal/tester"
	"golang.org/x/net/html"
)

func TestRemoveRedundant(t *testing.T) {
	runs := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:  "don't change other tags",
			input: `<span>a</span> <span>b</span>`,
			expected: `
├─body
│ ├─span
│ │ ├─#text "a"
│ ├─#text " "
│ ├─span
│ │ ├─#text "b"
			`,
		},
		{
			desc:  "don't change simple strong",
			input: `<strong>a</strong>`,

			expected: `
├─body
│ ├─strong
│ │ ├─#text "a"
			`,
		},
		{
			desc:  "remove double strong",
			input: `<strong><strong>a</strong></strong>`,

			expected: `
├─body
│ ├─strong
│ │ ├─#text "a"
			`,
		},
		{
			desc:  "remove more complicated double strong",
			input: `<strong><strong>a</strong> b <strong><strong>c</strong></strong></strong>`,

			expected: `
├─body
│ ├─strong
│ │ ├─#text "a"
│ │ ├─#text " b "
│ │ ├─#text "c"
			`,
		},

		{
			desc:  "leave italic inside bold",
			input: `<strong>A<em>B</em>C</strong>`,

			expected: `
├─body
│ ├─strong
│ │ ├─#text "A"
│ │ ├─em
│ │ │ ├─#text "B"
│ │ ├─#text "C"
			`,
		},
		{
			desc:  "dont leave other italic inside another italic",
			input: `<i>A<em>B</em>C</i>`,

			expected: `
├─body
│ ├─i
│ │ ├─#text "A"
│ │ ├─#text "B"
│ │ ├─#text "C"
			`,
		},
	}
	for _, run := range runs {
		t.Run(run.desc, func(t *testing.T) {
			doc := tester.Parse(t, run.input, "")

			RemoveRedundant(doc, func(a, b *html.Node) bool {
				isItalic := func(n *html.Node) bool {
					name := dom.NodeName(n)
					return name == "em" || name == "i"
				}
				isBold := func(n *html.Node) bool {
					name := dom.NodeName(n)
					return name == "strong" || name == "b"
				}

				if isItalic(a) && isItalic(b) {
					return true
				}
				if isBold(a) && isBold(b) {
					return true
				}

				return false
			})

			tester.ExpectRepresentation(t, doc, "output", run.expected)
		})
	}
}
