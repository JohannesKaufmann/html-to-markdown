package domutils

import (
	"context"
	"testing"

	"github.com/JohannesKaufmann/html-to-markdown/v2/internal/tester"
)

func TestRenameFakeSpans(t *testing.T) {
	runs := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:  "don't change other tags",
			input: `<p>a</p> <p>b</p>`,
			expected: `
├─body
│ ├─p
│ │ ├─#text "a"
│ ├─#text " "
│ ├─p
│ │ ├─#text "b"
			`,
		},
		{
			desc:  "don't change simple span",
			input: `<span>a</span>`,

			expected: `
├─body
│ ├─span
│ │ ├─#text "a"
			`,
		},
		{
			desc:  "don't change span with inline element",
			input: `<span><a>link content</a></span>`,

			expected: `
├─body
│ ├─span
│ │ ├─a
│ │ │ ├─#text "link content"
			`,
		},
		{
			desc:  "change span with block element",
			input: `<span><p>paragraph content</p></span>`,

			expected: `
├─body
│ ├─div
│ │ ├─p
│ │ │ ├─#text "paragraph content"
			`,
		},
		{
			desc:  "change multiple spans with block element",
			input: `<span><span><p>paragraph content</p></span></span>`,

			expected: `
├─body
│ ├─div
│ │ ├─div
│ │ │ ├─p
│ │ │ │ ├─#text "paragraph content"
			`,
		},
	}
	for _, run := range runs {
		t.Run(run.desc, func(t *testing.T) {
			doc := tester.Parse(t, run.input, "")

			RenameFakeSpans(context.TODO(), doc)

			tester.ExpectRepresentation(t, doc, "output", run.expected)
		})
	}
}
