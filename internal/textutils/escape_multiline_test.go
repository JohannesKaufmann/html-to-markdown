package textutils

import (
	"bytes"
	"strings"
	"testing"
)

func EscapeMultiLine_Old(content []byte) []byte {
	content = bytes.TrimSpace(content)
	content = TrimConsecutiveNewlines(content)
	if len(content) == 0 {
		return content
	}

	parts := bytes.Split(content, newline)
	for i := range parts {
		parts[i] = bytes.TrimSpace(parts[i])
		if len(parts[i]) == 0 {
			parts[i] = escape
		}
	}
	content = bytes.Join(parts, newline)

	return content
}

func TestEscapeMultiLine(t *testing.T) {
	var tests = []struct {
		Name     string
		Text     string
		Expected string
	}{
		{
			Name:     "empty",
			Text:     "",
			Expected: "",
		},
		{
			Name:     "not needed",
			Text:     "some longer text that is on one line",
			Expected: "some longer text that is on one line",
		},

		{
			Name:     "one newline",
			Text:     "A\nB",
			Expected: "A\nB",
		},
		{
			Name:     "two newlines",
			Text:     "A\n\nB",
			Expected: "A\n\\\nB",
		},
		{

			Name: "many newlines",
			// Will be max two newlines characters
			Text:     "line 1\n\n\n\nline 2",
			Expected: "line 1\n\\\nline 2",
		},

		{
			Name: "multiple empty lines",
			Text: `line1
line2

line3




line4`,
			Expected: `line1
line2
\
line3
\
line4`,
		},

		{
			Name:     "empty line with a space",
			Text:     "line 1\n  \nline 2",
			Expected: "line 1\n\\\nline 2",
		},

		{
			Name:     "content has a space",
			Text:     "a\n\n b",
			Expected: "a\n\\\nb",
		},
		{
			Name:     "content is indented",
			Text:     "line 1\n  line 2\n\tline 3",
			Expected: "line 1\nline 2\nline 3",
		},

		// TODO: keep existing "\" characters?
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Run("old", func(t *testing.T) {
				output := EscapeMultiLine_Old([]byte(test.Text))

				if string(output) != test.Expected {
					t.Errorf("expected '%s' but got '%s'", test.Expected, string(output))
				}
			})
			t.Run("new", func(t *testing.T) {
				output := EscapeMultiLine([]byte(test.Text))

				if string(output) != test.Expected {
					t.Errorf("expected '%s' but got '%s'", test.Expected, string(output))
				}
			})

		})

	}
}

func BenchmarkEscapeMultiLine(b *testing.B) {

	b.Run("old", func(b *testing.B) {
		input := []byte(strings.Repeat("line 1\n\n  \nline 2", 100))

		for i := 0; i < b.N; i++ {
			_ = EscapeMultiLine_Old(input)
		}
	})
	b.Run("new", func(b *testing.B) {
		input := []byte(strings.Repeat("line 1\n\n  \nline 2", 100))

		for i := 0; i < b.N; i++ {
			_ = EscapeMultiLine(input)
		}
	})
}
