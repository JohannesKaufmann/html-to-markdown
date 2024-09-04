package converter

import (
	"testing"

	"golang.org/x/net/html"
)

func TestEscapeContent(t *testing.T) {
	conv := NewConverter()
	conv.Register.EscapedChar('>')

	input := []byte{'a', '>'}

	output := conv.escapeContent(input)
	if len(output) != 3 {
		t.Error("expected different length")
	}
	// Since '>' is a character used for quotes in markdown,
	// there should be a marker before it.
	if output[0] != 'a' || output[1] != placeholderByte || output[2] != '>' {
		t.Error("expected different characters")
	}
}

func TestUnEscapeContent(t *testing.T) {
	conv := NewConverter()
	conv.Register.UnEscaper(func(chars []byte, index int) int {
		if chars[index] != '>' {
			return -1
		}

		// A bit too simplistic for demonstration purposes.
		// Normally here would be content to check if the escaping is necessary...
		return 1
	}, PriorityStandard)

	input := []byte{placeholderByte, 'a', 'b'}
	output := conv.unEscapeContent(input)
	if len(output) != 2 {
		t.Error("expected different length")
	}
	// No escaping needed
	if output[0] != 'a' || output[1] != 'b' {
		t.Error("expected different characters")
	}

	input = []byte{placeholderByte, '>', 'a'}
	output = conv.unEscapeContent(input)
	if len(output) != 3 {
		t.Error("expected different length")
	}
	// Escaping needed
	if output[0] != '\\' || output[1] != '>' || output[2] != 'a' {
		t.Error("expected different characters")
	}
}

func TestWithEscapeMode(t *testing.T) {
	testCases := []struct {
		desc string
		mode escapeMode

		input    string
		expected string
	}{
		{
			desc: "EscapeSmart",
			mode: EscapeSmart,

			input:    "a|b",
			expected: "a\\|b",
		},
		{
			desc: "EscapeDisabled",
			mode: EscapeDisabled,

			input:    "a|b",
			expected: "a|b",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			conv := NewConverter(
				WithEscapeMode(tC.mode),
			)
			conv.Register.Renderer(func(ctx Context, w Writer, n *html.Node) RenderStatus {
				return RenderTryNext
			}, PriorityStandard)

			conv.Register.EscapedChar('|')
			conv.Register.UnEscaper(func(chars []byte, index int) int {
				if chars[index] != '|' {
					return -1
				}

				// A bit too simplistic for demonstration purposes.
				// Normally here would be content to check if the escaping is necessary...
				return 1
			}, PriorityStandard)

			output, err := conv.ConvertString(tC.input)
			if err != nil {
				t.Error(err)
			}

			if output != tC.expected {
				t.Errorf("expected %q but got %q", tC.expected, output)
			}
		})
	}
}
