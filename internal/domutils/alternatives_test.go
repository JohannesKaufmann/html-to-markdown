package domutils

import (
	"context"
	"testing"

	"github.com/JohannesKaufmann/html-to-markdown/v2/internal/tester"
)

func TestLeafBlockAlternatives(t *testing.T) {
	runs := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:  "divider in heading",
			input: `<h3>Heading<hr /></h3>`,
			expected: `
├─body
│ ├─h3
│ │ ├─#text "Heading"
			`,
		},
		{
			desc:  "simple",
			input: `<a href="/page.html"><h3>Heading</h3></a>`,
			expected: `
├─body
│ ├─a (href="/page.html")
│ │ ├─strong
│ │ │ ├─#text "Heading"
│ │ ├─br
			`,
		},
		{
			desc:  "two headings",
			input: `<a href="/page.html"><h4>Heading A</h4><h3>Heading B</h3></a>`,
			expected: `
├─body
│ ├─a (href="/page.html")
│ │ ├─strong
│ │ │ ├─#text "Heading A"
│ │ ├─br
│ │ ├─strong
│ │ │ ├─#text "Heading B"
│ │ ├─br
			`,
		},
		{
			desc: "two headings formatted",
			input: `
<a href="/page.html">
	<h4>Heading A</h4>
	<h3>Heading B</h3>
</a>
			`,
			expected: `
├─body
│ ├─a (href="/page.html")
│ │ ├─#text "\n\t"
│ │ ├─strong
│ │ │ ├─#text "Heading A"
│ │ ├─br
│ │ ├─#text "\n\t"
│ │ ├─strong
│ │ │ ├─#text "Heading B"
│ │ ├─br
│ │ ├─#text "\n"
			`,
		},
	}
	for _, run := range runs {
		t.Run(run.desc, func(t *testing.T) {
			doc := tester.Parse(t, run.input, "")

			LeafBlockAlternatives(context.TODO(), doc)

			tester.ExpectRepresentation(t, doc, "output", run.expected)
		})
	}
}
