package md

import (
	"fmt"
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
