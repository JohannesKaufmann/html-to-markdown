package escape

import (
	"reflect"
	"testing"
)

func TestIsDivider(t *testing.T) {
	runs := []struct {
		name  string
		chars []byte

		expected []int
	}{
		{
			name:     "not needed",
			chars:    []byte{'a', 'b', 'c'},
			expected: []int{-1, -1, -1},
		},
		{
			name:     "two dashes",
			chars:    []byte{'-', '-'},
			expected: []int{-1, -1},
		},
		{
			name:     "char afterwards",
			chars:    []byte{'-', '-', '-', ' ', 'a'},
			expected: []int{-1, -1, -1, -1, -1},
		},

		{
			name:     "three dashes",
			chars:    []byte{'-', '-', '-'},
			expected: []int{3, -1, -1},
		},
		{
			name:     "five dashes",
			chars:    []byte{'-', '-', '-', '-', '-'},
			expected: []int{5, -1, -1, -1, -1},
		},
		{
			name:     "space after",
			chars:    []byte{'-', '-', '-', ' '},
			expected: []int{4, -1, -1, -1},
		},
		{
			name:     "newline after",
			chars:    []byte{'-', '-', '-', '\n'},
			expected: []int{3, -1, -1, -1},
		},

		{
			name:     "newline and space before",
			chars:    []byte{'\n', ' ', '-', '-', '-'},
			expected: []int{-1, -1, 3, -1, -1},
		},
		{
			name:     "with placeholders",
			chars:    []byte{'\n', ' ', placeholderByte, '-', placeholderByte, '-', placeholderByte, '-'},
			expected: []int{-1, -1, -1, 5, -1, -1, -1, -1},
		},
	}
	for _, run := range runs {

		t.Run(run.name, func(t *testing.T) {
			var actual []int
			for index := range run.chars {
				output := IsDivider(run.chars, index)

				actual = append(actual, output)
			}

			if !reflect.DeepEqual(actual, run.expected) {
				t.Errorf("expected %+v but got %+v", run.expected, actual)
			}
		})

	}
}
