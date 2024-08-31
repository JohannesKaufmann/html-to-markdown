package textutils

import "testing"

func TestCollapseInlineCodeContent(t *testing.T) {
	runs := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:     "empty",
			input:    "",
			expected: "",
		},
		{
			desc:     "not needed",
			input:    "a b",
			expected: "a b",
		},
		{
			desc:     "one newline",
			input:    "a\nb",
			expected: "a b",
		},
		{
			desc:     "multiple newlines",
			input:    "a\nb\n\nc",
			expected: "a b c",
		},
		{
			desc:     "also trim",
			input:    " a b ",
			expected: "a b",
		},
		{
			desc: "realistic code content",
			input: `
			
			body {
				color: yellow;
				font-size: 16px;
			}
			
			`,
			expected: "body { color: yellow; font-size: 16px; }",
		},
	}
	for _, run := range runs {
		t.Run(run.desc, func(t *testing.T) {
			actual := CollapseInlineCodeContent([]byte(run.input))
			if string(actual) != run.expected {
				t.Errorf("expected %q but got %q", run.expected, string(actual))
			}
		})
	}
}
