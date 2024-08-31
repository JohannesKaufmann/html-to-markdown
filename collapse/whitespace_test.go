package collapse

import (
	"regexp"
	"strings"
	"testing"
)

// This is the alternative (but slower) function that uses regex.
func _regexReplaceAnyWhitespaceWithSpace(text string) string {
	var rAnyWhitespace = regexp.MustCompile(`[ \r\n\t]+`)

	return rAnyWhitespace.ReplaceAllString(text, " ")
}

func TestReplaceAnyWhitespaceWithSpace(t *testing.T) {
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
			desc:     "one space",
			input:    " ",
			expected: " ",
		},
		{
			desc:     "two spaces",
			input:    "  ",
			expected: " ",
		},
		{
			desc:     "many spaces",
			input:    "     ",
			expected: " ",
		},
		{
			desc:     "one newline",
			input:    "\n",
			expected: " ",
		},
		{
			desc:     "many newlines",
			input:    "\n\n\n\n",
			expected: " ",
		},
		{
			desc:     "combination of newlines and spaces",
			input:    "\n a \nb  \nc\n",
			expected: " a b c ",
		},
		{
			desc:     "special dash",
			input:    " \u2013 ",
			expected: " \u2013 ",
		},
		{
			desc:     "no spaces in text",
			input:    "abcdef",
			expected: "abcdef",
		},
		{
			desc:     "one space in text",
			input:    "abc def",
			expected: "abc def",
		},
		{
			desc:     "two spaces in text",
			input:    "abc  def",
			expected: "abc def",
		},
		{
			desc:     "one newline in text",
			input:    "abc\ndef",
			expected: "abc def",
		},
		{
			desc:     "two newlines in text",
			input:    "abc\n\ndef",
			expected: "abc def",
		},
		{
			desc:     "a newline and space in text",
			input:    "abc \ndef",
			expected: "abc def",
		},
		{
			desc:     "one space before text",
			input:    " abcdef",
			expected: " abcdef",
		},
		{
			desc:     "two spaces before text",
			input:    "  abcdef",
			expected: " abcdef",
		},
		{
			desc:     "one space after text",
			input:    "abcdef ",
			expected: "abcdef ",
		},
		{
			desc:     "two spaces after text",
			input:    "abcdef  ",
			expected: "abcdef ",
		},
		{
			desc:     "multiple spaces before & one space after",
			input:    "   or ",
			expected: " or ",
		},
		{
			desc:     "multiple spaces before & multiple spaces after",
			input:    "   or   ",
			expected: " or ",
		},
		{
			desc:     "one space before & multiple spaces after",
			input:    " or   ",
			expected: " or ",
		},
	}

	for _, run := range runs {
		t.Run(run.desc, func(t *testing.T) {
			t.Run("Regex Version", func(t *testing.T) {
				output := _regexReplaceAnyWhitespaceWithSpace(run.input)

				if output != run.expected {
					t.Errorf("expected %q but got %q", run.expected, output)
				}
			})

			t.Run("New Version", func(t *testing.T) {
				output := replaceAnyWhitespaceWithSpace(run.input)
				if output != run.expected {
					t.Errorf("expected %q but got %q", run.expected, output)
				}

				// Instead of writing all tests twice...
				output2 := replaceAnyWhitespaceWithSpace(strings.ReplaceAll(run.input, " ", "\n"))
				if output2 != run.expected {
					t.Errorf("for newlines: expected %q but got %q", run.expected, output2)
				}
			})

		})
	}
}

func FuzzReplaceAnyWhitespaceWithSpace(f *testing.F) {
	f.Add("abc def")
	f.Add(" ")
	f.Add("abc\n\ndef")

	f.Fuzz(func(t *testing.T, orig string) {
		output1 := _regexReplaceAnyWhitespaceWithSpace(orig)
		output2 := replaceAnyWhitespaceWithSpace(orig)

		if output1 != output2 {
			t.Errorf("input:%q => regex: %q function: %q", orig, output1, output2)
		}
	})
}

func TestReplaceAnyWhitespaceWithSpace_Allocs(t *testing.T) {
	const N = 1000

	runs := []struct {
		desc           string
		input          string
		expectedAllocs float64
	}{
		{
			desc:           "empty string",
			input:          "",
			expectedAllocs: 0,
		},
		{
			desc:           "one space",
			input:          " ",
			expectedAllocs: 0,
		},
		{
			desc:           "no spaces",
			input:          "abcdef",
			expectedAllocs: 0,
		},
		{
			desc:           "one space at start",
			input:          " abcdef",
			expectedAllocs: 0,
		},
		{
			desc:           "one space at end",
			input:          "abcdef ",
			expectedAllocs: 0,
		},
		{
			desc:           "one space in middle",
			input:          "abc def",
			expectedAllocs: 0,
		},
		{
			desc:           "one space at start, middle and end",
			input:          " abc def ",
			expectedAllocs: 0,
		},
		{
			desc:           "one space at start, middle and end",
			input:          "some longer text with spaces",
			expectedAllocs: 0,
		},

		{
			desc:           "multiple spaces",
			input:          "abc    def",
			expectedAllocs: 1,
		},
		{
			desc:           "multiple newlines & spaces",
			input:          "\n\nab  cdef  \n",
			expectedAllocs: 1,
		},
		{
			desc:           "longer string",
			input:          strings.Repeat("Lorem Ipsum is simply dummy text", 10),
			expectedAllocs: 0,
		},
	}
	for _, run := range runs {
		t.Run(run.desc, func(t *testing.T) {
			avg := testing.AllocsPerRun(N, func() {
				output := replaceAnyWhitespaceWithSpace(run.input)
				_ = output
			})
			if avg != run.expectedAllocs {
				t.Errorf("expected %f allocations but got %f", run.expectedAllocs, avg)
			}
		})
	}
}
