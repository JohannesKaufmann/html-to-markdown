package escape

import (
	"reflect"
	"testing"
)

func TestIsBlockquote(t *testing.T) {
	runs := []struct {
		name  string
		chars []byte

		expected []int
	}{
		{
			name:     "allow simple quote",
			chars:    []byte{'>', ' ', 'a'},
			expected: []int{1, -1, -1},
		},
		{
			name:     "allow space before",
			chars:    []byte{' ', '>', ' ', 'a'},
			expected: []int{-1, 1, -1, -1},
		},
		{
			name:     "allow missing space after",
			chars:    []byte{'>', 'a'},
			expected: []int{1, -1},
		},
		{
			name:     "allow newline before",
			chars:    []byte{'\n', '>', ' ', 'a'},
			expected: []int{-1, 1, -1, -1},
		},
		{
			name:     "allow newline and space before",
			chars:    []byte{'\n', ' ', '>', ' ', 'a'},
			expected: []int{-1, -1, 1, -1, -1},
		},
		{
			name:     "allow placeholder before",
			chars:    []byte{placeholderByte, '>', ' ', 'a'},
			expected: []int{-1, 1, -1, -1},
		},
		{
			name:     "dont allow other chars before",
			chars:    []byte{'a', '>', ' ', 'a'},
			expected: []int{-1, -1, -1, -1},
		},
	}
	for _, run := range runs {

		t.Run(run.name, func(t *testing.T) {
			var actual []int
			for index := range run.chars {
				output := IsBlockQuote(run.chars, index)

				actual = append(actual, output)
			}

			if !reflect.DeepEqual(actual, run.expected) {
				t.Errorf("expected %+v but got %+v", run.expected, actual)
			}
		})

	}
}
