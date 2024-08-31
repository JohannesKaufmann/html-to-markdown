package textutils

import (
	"strings"
	"testing"
)

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
