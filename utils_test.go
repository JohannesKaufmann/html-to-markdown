package md

import (
	"fmt"
	"strings"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func TestTrimpLeadingSpaces(t *testing.T) {
	var tests = []struct {
		Name   string
		Text   string
		Expect string
	}{
		{
			Name: "trim normal text",
			Text: `
This is a normal paragraph
 this as well
  just with some spaces before
			`,
			Expect: `
This is a normal paragraph
this as well
just with some spaces before
			`,
		},
		{
			Name: "dont trim nested lists",
			Text: `
- Home
- About
	- People
	- History
		- 2019
		- 2020		
			`,
			Expect: `
- Home
- About
	- People
	- History
		- 2019
		- 2020		
			`,
		},
		{
			Name: "dont trim list with multiple paragraphs",
			Text: `
1.  This is a list item with two paragraphs. Lorem ipsum dolor
	sit amet, consectetuer adipiscing elit. Aliquam hendrerit
	mi posuere lectus.

	Vestibulum enim wisi, viverra nec, fringilla in, laoreet
	vitae, risus. Donec sit amet nisl. Aliquam semper ipsum
	sit amet velit.

2.  Suspendisse id sem consectetuer libero luctus adipiscing.
			`,
			Expect: `
1.  This is a list item with two paragraphs. Lorem ipsum dolor
	sit amet, consectetuer adipiscing elit. Aliquam hendrerit
	mi posuere lectus.

	Vestibulum enim wisi, viverra nec, fringilla in, laoreet
	vitae, risus. Donec sit amet nisl. Aliquam semper ipsum
	sit amet velit.

2.  Suspendisse id sem consectetuer libero luctus adipiscing.
			`,
		},
		{
			Name: "dont trim code blocks",
			Text: `
This is a normal paragraph:

    This is a code block.
			`,
			Expect: `
This is a normal paragraph:

    This is a code block.
			`,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			output := TrimpLeadingSpaces(test.Text)

			if output != test.Expect {
				dmp := diffmatchpatch.New()
				diffs := dmp.DiffMain(test.Expect, output, false)

				fmt.Println(dmp.DiffToDelta(diffs))
				fmt.Println(dmp.DiffPrettyText(diffs))

				t.Errorf("expected '%s' but got '%s'", test.Expect, output)
			}
		})
	}

}

func TestTrimTrailingSpaces(t *testing.T) {
	var tests = []struct {
		Name   string
		Text   string
		Expect string
	}{
		{
			Name: "trim after normal text",
			Text: `
1\. xxx 

2\. xxxx	
			`,
			Expect: `
1\. xxx

2\. xxxx
`,
		},
		{
			Name:   "dont trim inside normal text",
			Text:   "When `x = 3`, that means `x + 2 = 5`",
			Expect: "When `x = 3`, that means `x + 2 = 5`",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			output := TrimTrailingSpaces(test.Text)

			if output != test.Expect {
				dmp := diffmatchpatch.New()
				diffs := dmp.DiffMain(test.Expect, output, false)

				fmt.Println(dmp.DiffToDelta(diffs))
				fmt.Println(dmp.DiffPrettyText(diffs))

				t.Errorf("expected '%s' but got '%s'", test.Expect, output)
			}
		})
	}
}

func TestEscapeMultiLine(t *testing.T) {
	var tests = []struct {
		Name   string
		Text   string
		Expect string
	}{
		{
			Name: "escape new lines",
			Text: `line1
line2

line3









line4`,
			Expect: `line1\
line2\
\
\
line3\
\
\
line4`,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			output := EscapeMultiLine(test.Text)

			if output != test.Expect {
				dmp := diffmatchpatch.New()
				diffs := dmp.DiffMain(test.Expect, output, false)

				fmt.Println(dmp.DiffToDelta(diffs))
				fmt.Println(dmp.DiffPrettyText(diffs))

				t.Errorf("expected '%s' but got '%s'", test.Expect, output)
			}
		})
	}
}

func TestCalculateCodeFence(t *testing.T) {
	var tests = []struct {
		Name      string
		FenceChar rune

		Text   string
		Expect string
	}{
		{
			Name:      "no occurrences with backtick",
			FenceChar: '`',
			Text:      `normal ~~~ code block`,
			Expect:    "```",
		},
		{
			Name:      "no occurrences with tilde",
			FenceChar: '~',
			Text:      "normal ``` code block",
			Expect:    "~~~",
		},
		{
			Name:      "one exact occurrence",
			FenceChar: '`',
			Text:      "```",
			Expect:    "````",
		},
		{
			Name:      "one occurrences with backtick",
			FenceChar: '`',
			Text:      "normal ``` code block",
			Expect:    "````",
		},
		{
			Name:      "one bigger occurrences with backtick",
			FenceChar: '`',
			Text:      "normal ````` code block",
			Expect:    "``````",
		},
		{
			Name:      "multiple occurrences with backtick",
			FenceChar: '`',
			Text:      "normal ``` code `````` block",
			Expect:    "```````",
		},
		{
			Name:      "multiple occurrences with tilde",
			FenceChar: '~',
			Text:      "normal ~~~ code ~~~~~~~~~~~~ block",
			Expect:    "~~~~~~~~~~~~~",
		},
		{
			Name:      "multiple occurrences on different lines with tilde",
			FenceChar: '~',
			Text: `
normal
	~~~
code ~~~~~~~~~~~~ block
				`,
			Expect: "~~~~~~~~~~~~~",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			output := CalculateCodeFence(test.FenceChar, test.Text)

			if output != test.Expect {
				t.Errorf("expected '%s' (x%d) but got '%s' (x%d)", test.Expect, strings.Count(test.Expect, string(test.FenceChar)), output, strings.Count(output, string(test.FenceChar)))
			}
		})
	}
}
