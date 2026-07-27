package domutils

import (
	"context"
	"testing"

	"github.com/JohannesKaufmann/dom"
	"github.com/JohannesKaufmann/html-to-markdown/v2/internal/tester"
	"golang.org/x/net/html"
)

// The runs are the examples that the commonmark spec itself lists in
// section 6.2, expressed as the characters surrounding the delimiter run.
func TestFlanking(t *testing.T) {
	runs := []struct {
		desc     string
		prev     rune
		next     rune
		canOpen  bool
		canClose bool
	}{
		{desc: `***abc`, prev: 0, next: 'a', canOpen: true, canClose: false},
		{desc: `**"abc"`, prev: 0, next: '"', canOpen: true, canClose: false},
		{desc: ` abc***`, prev: 'c', next: 0, canOpen: false, canClose: true},
		{desc: `"abc"**`, prev: '"', next: 0, canOpen: false, canClose: true},
		{desc: ` abc***def`, prev: 'c', next: 'd', canOpen: true, canClose: true},
		{desc: `"abc"_"def"`, prev: '"', next: '"', canOpen: true, canClose: true},
		{desc: `abc *** def`, prev: ' ', next: ' ', canOpen: false, canClose: false},
		{desc: `a _ b`, prev: ' ', next: ' ', canOpen: false, canClose: false},

		// The cases that this package has to move punctuation for.
		{desc: `before*.middle`, prev: 'e', next: '.', canOpen: false, canClose: true},
		{desc: `middle.*after`, prev: '.', next: 'a', canOpen: true, canClose: false},
		{desc: `before.*middle`, prev: '.', next: 'm', canOpen: true, canClose: false},
		{desc: `middle*.after`, prev: 'e', next: '.', canOpen: false, canClose: true},

		// A currency sign is in the S category, which commonmark counts as punctuation.
		{desc: `before*€x`, prev: 'e', next: '€', canOpen: false, canClose: true},
		// A CJK ideograph is a letter, so it does not.
		{desc: `before*日x`, prev: 'e', next: '日', canOpen: true, canClose: true},
	}
	for _, run := range runs {
		t.Run(run.desc, func(t *testing.T) {
			if got := CanOpen(run.prev, run.next); got != run.canOpen {
				t.Errorf("CanOpen(%q, %q) = %v but want %v", run.prev, run.next, got, run.canOpen)
			}
			if got := CanClose(run.prev, run.next); got != run.canClose {
				t.Errorf("CanClose(%q, %q) = %v but want %v", run.prev, run.next, got, run.canClose)
			}
		})
	}
}

func TestMovePunctuation(t *testing.T) {
	runs := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:  "punctuation at the start",
			input: `before<em>.middle</em>after`,
			expected: `
├─body
│ ├─#text "before."
│ ├─em
│ │ ├─#text "middle"
│ ├─#text "after"
			`,
		},
		{
			desc:  "punctuation at the end",
			input: `before<em>middle.</em>after`,
			expected: `
├─body
│ ├─#text "before"
│ ├─em
│ │ ├─#text "middle"
│ ├─#text ".after"
			`,
		},
		{
			desc:  "one character is enough",
			input: `before<em>?!!middle?!!</em>after`,
			expected: `
├─body
│ ├─#text "before?"
│ ├─em
│ │ ├─#text "!!middle?!"
│ ├─#text "!after"
			`,
		},
		{
			desc:  "surrounded by spaces, so nothing to do",
			input: `before <em>.middle.</em> after`,
			expected: `
├─body
│ ├─#text "before "
│ ├─em
│ │ ├─#text ".middle."
│ ├─#text " after"
			`,
		},
		{
			desc:  "surrounded by punctuation, so nothing to do",
			input: `before-<em>.middle.</em>-after`,
			expected: `
├─body
│ ├─#text "before-"
│ ├─em
│ │ ├─#text ".middle."
│ ├─#text "-after"
			`,
		},
		{
			desc:  "no text to move it into",
			input: `<em>.middle.</em>`,
			expected: `
├─body
│ ├─em
│ │ ├─#text ".middle."
			`,
		},
		{
			desc:  "a span has no markdown meaning",
			input: `<span>before</span><em>.middle</em><span>after</span>`,
			expected: `
├─body
│ ├─span
│ │ ├─#text "before."
│ ├─em
│ │ ├─#text "middle"
│ ├─span
│ │ ├─#text "after"
			`,
		},
		{
			desc:  "the content is only punctuation",
			input: `before<em>.</em>after`,
			expected: `
├─body
│ ├─#text "before."
│ ├─em
│ │ ├─#text ""
│ ├─#text "after"
			`,
		},
	}
	for _, run := range runs {
		t.Run(run.desc, func(t *testing.T) {
			doc := tester.Parse(t, run.input, "")

			MovePunctuation(context.Background(), doc, func(n *html.Node) bool {
				return dom.NodeName(n) == "em"
			})

			tester.ExpectRepresentation(t, doc, "output", run.expected)
		})
	}
}
