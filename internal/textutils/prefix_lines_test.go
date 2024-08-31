package textutils

import (
	"bytes"
	"regexp"
	"testing"
)

var beginningR = regexp.MustCompile(`(?m)^`)

func _oldPrefixLines(content string, repl string) string {
	return beginningR.ReplaceAllString(content, repl)
}

func TestPrefixLines(t *testing.T) {
	runs := []struct {
		desc     string
		input    []byte
		expected []byte
	}{
		{
			desc:     "one line",
			input:    []byte("abc"),
			expected: []byte("> abc"),
		},
		{
			desc:     "two lines",
			input:    []byte("line 1\nline 2"),
			expected: []byte("> line 1\n> line 2"),
		},
		{
			desc:     "two newlines between",
			input:    []byte("line 1\n\nline 2"),
			expected: []byte("> line 1\n> \n> line 2"),
		},
		{
			desc:     "newline at end",
			input:    []byte("abc\n"),
			expected: []byte("> abc\n> "),
		},
	}

	for _, run := range runs {
		t.Run(run.desc, func(t *testing.T) {
			t.Run("old", func(t *testing.T) {
				output := _oldPrefixLines(string(run.input), "> ")

				if output != string(run.expected) {
					t.Errorf("expected %q but got %q", string(run.expected), output)
				}
			})
			t.Run("new", func(t *testing.T) {
				output := PrefixLines(run.input, []byte{'>', ' '})

				if !bytes.Equal(output, run.expected) {
					t.Errorf("expected %q but got %q", string(run.expected), string(output))
				}
			})
		})
	}
}
