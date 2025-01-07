package textutils

import (
	"bytes"
	"testing"
)

func TestTrimConsecutiveNewlines(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"single char", "a", "a"},
		{"simple text", "hello", "hello"},
		{"normal text without newlines", "hello  this is a   normal text", "hello  this is a   normal text"},

		// Single newline cases
		{"single newline", "a\nb", "a\nb"},
		{"single newline with spaces", "a  \nb", "a  \nb"},
		{"spaces after newline", "a\n  b", "a\n  b"},

		// Double newline cases
		{"double newline", "a\n\nb", "a\n\nb"},
		{"double newline with spaces", "a  \n\nb", "a\n\nb"},
		{"spaces between newlines", "a\n  \nb", "a\n  \nb"},
		// Note: It should not change the spaces *after* the newlines since this could impacts lists
		{"spaces after double newline", "a\n\n  b", "a\n\n  b"},

		// Triple+ newline cases
		{"triple newline", "a\n\n\nb", "a\n\nb"},
		{"quad newline", "a\n\n\n\nb", "a\n\nb"},
		{"triple newline with spaces", "a  \n\n\nb", "a\n\nb"},

		// Multiple segment cases
		{"multiple segments", "a\n\nb\n\nc", "a\n\nb\n\nc"},
		{"multiple segments with spaces", "a  \n\nb  \n\nc", "a\n\nb\n\nc"},

		// Spaces at end of line
		{"hard-line-break followed by text", "a  \nb", "a  \nb"},
		{"hard-line-break followed by newline", "a  \n\nb", "a\n\nb"},

		// Edge cases
		{"only newlines", "\n\n\n", "\n\n"},
		{"only spaces", "   ", "   "},

		{"leading and trailing newlines", "\n\n\ntext\n\n\n", "\n\ntext\n\n"},
		{"newlines and spaces 1", "  \n  \n  \n  \n  ", "\n\n  "},
		{"newlines and spaces 2", "a  \n  \nb", "a\n\nb"},
		{"newlines and spaces 3", "a  \n \nb", "a\n\nb"},

		{"leading spaces", "   a", "   a"},
		{"leading newline 1", "\na", "\na"},
		{"leading newline 2", "\n\na", "\n\na"},
		{"leading newline 3", "\n\n\na", "\n\na"},

		{"trailing spaces", "a   ", "a   "},
		{"trailing newline 1", "a\n", "a\n"},
		{"trailing newlines 2", "a\n\n", "a\n\n"},
		{"trailing newlines 3", "a\n\n\n", "a\n\n"},

		// UTF-8 cases
		{"german special chars", "äöü\n\n\näöü", "äöü\n\näöü"},
		{"utf8 chars", "🌟\n\n\n🌟\n\n\n🌟", "🌟\n\n🌟\n\n🌟"},

		// Markdown
		// Note: The sublist needs to be indented by "  -"
		{"indented sublist", "- The main list\n  \n  - The sublist", "- The main list\n  \n  - The sublist"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := TrimConsecutiveNewlines([]byte(tt.input))
			output = TrimUnnecessaryHardLineBreaks(output)

			if string(output) != tt.expected {
				t.Errorf("\ninput:    %q\nexpected: %q\ngot:      %q",
					tt.input, tt.expected, string(output),
				)
			}
		})
	}
}

func TestTrimConsecutiveNewlines_Allocs(t *testing.T) {
	const N = 1000

	t.Run("no newlines", func(t *testing.T) {
		var expectedAverage float64 = 1

		actualAverage := testing.AllocsPerRun(N, func() {
			input := []byte("abc")
			output := TrimConsecutiveNewlines(input)
			_ = output
		})
		if actualAverage != expectedAverage {
			t.Errorf("expected %f allocations but got %f", expectedAverage, actualAverage)
		}
	})
	t.Run("exactly two newlines", func(t *testing.T) {
		var expectedAverage float64 = 1

		actualAverage := testing.AllocsPerRun(N, func() {
			input := []byte("abc\n\nabc")
			output := TrimConsecutiveNewlines(input)
			_ = output
		})
		if actualAverage != expectedAverage {
			t.Errorf("expected %f allocations but got %f", expectedAverage, actualAverage)
		}
	})
	t.Run("three newlines", func(t *testing.T) {
		var expectedAverage float64 = 1

		actualAverage := testing.AllocsPerRun(N, func() {
			input := []byte("abc\n\n\nabc")
			output := TrimConsecutiveNewlines(input)
			_ = output
		})
		if actualAverage != expectedAverage {
			t.Errorf("expected %f allocations but got %f", expectedAverage, actualAverage)
		}
	})
	t.Run("many newlines", func(t *testing.T) {
		var expectedAverage float64 = 16

		actualAverage := testing.AllocsPerRun(N, func() {
			input := bytes.Repeat([]byte("abc\n\n\n\n\n\nabc"), 1000)
			output := TrimConsecutiveNewlines(input)
			_ = output
		})
		if actualAverage != expectedAverage {
			t.Errorf("expected %f allocations but got %f", expectedAverage, actualAverage)
		}
	})
}

func BenchmarkTrimConsecutiveNewlines(b *testing.B) {
	const Repeat = 10

	runs := []struct {
		desc  string
		input []byte
	}{
		{
			desc:  "not needed",
			input: bytes.Repeat([]byte("normal\n\ntext"), Repeat),
		},
		{
			desc:  "multiple times",
			input: bytes.Repeat([]byte("1\n\n\n2\n\n\n3"), Repeat),
		},
	}

	for _, run := range runs {
		b.Run(run.desc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				TrimConsecutiveNewlines(run.input)
			}
		})
	}
}
