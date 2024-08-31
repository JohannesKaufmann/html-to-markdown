package textutils

import (
	"bytes"
	"testing"
)

func TestSurroundingSpaces(t *testing.T) {
	testCases := []struct {
		desc  string
		input []byte

		expectedLeft    []byte
		expectedTrimmed []byte
		expectedRight   []byte
	}{
		{
			desc:  "empty string",
			input: []byte(""),

			expectedLeft:    []byte(""),
			expectedTrimmed: []byte(""),
			expectedRight:   []byte(""),
		},
		{
			desc:  "one space",
			input: []byte(" "),

			expectedLeft:    []byte(""),
			expectedTrimmed: []byte(""),
			expectedRight:   []byte(" "),
		},
		{
			desc:  "simple string",
			input: []byte("some text"),

			expectedLeft:    []byte(""),
			expectedTrimmed: []byte("some text"),
			expectedRight:   []byte(""),
		},
		{
			desc:  "spaces around",
			input: []byte("  text    "),

			expectedLeft:    []byte("  "),
			expectedTrimmed: []byte("text"),
			expectedRight:   []byte("    "),
		},
		{
			desc:  "newlines around",
			input: []byte("\n\n text  \n\n"),

			expectedLeft:    []byte("\n\n "),
			expectedTrimmed: []byte("text"),
			expectedRight:   []byte("  \n\n"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			leftExtra, trimmed, rightExtra := SurroundingSpaces(tC.input)

			if !bytes.Equal(leftExtra, tC.expectedLeft) {
				t.Errorf("expected %q but got %q for the left extra", string(tC.expectedLeft), string(leftExtra))
			}
			if !bytes.Equal(trimmed, tC.expectedTrimmed) {
				t.Errorf("expected %q but got %q for the trimmed text", string(tC.expectedTrimmed), string(trimmed))
			}
			if !bytes.Equal(rightExtra, tC.expectedRight) {
				t.Errorf("expected %q but got %q for the right extra", string(tC.expectedRight), string(rightExtra))
			}
		})
	}
}
