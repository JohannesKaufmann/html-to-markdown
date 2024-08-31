package escape

import (
	"reflect"
	"testing"
)

func TestIsAtxHeader(t *testing.T) {

	runs := []struct {
		name  string
		chars []rune

		expected []int
	}{
		{
			name:     "not needed",
			chars:    []rune{'a', 'b', 'c'},
			expected: []int{-1, -1, -1},
		},
		{
			name:     "inside text",
			chars:    []rune{'a', placeholderRune, '#', ' ', 'a'},
			expected: []int{-1, -1, -1, -1, -1},
		},
		{
			name:     "inside text with space between",
			chars:    []rune{'a', ' ', placeholderRune, '#', ' ', 'a'},
			expected: []int{-1, -1, -1, -1, -1, -1},
		},
		{
			name:     "h1 at start of file",
			chars:    []rune{placeholderRune, '#', ' ', 'a', 'b'},
			expected: []int{-1, 1, -1, -1, -1},
		},
		{
			name:     "h1 at start of line",
			chars:    []rune{'\n', placeholderRune, '#', ' ', 'a', 'b'},
			expected: []int{-1, -1, 1, -1, -1, -1},
		},
		{
			name:     "h1 with space before",
			chars:    []rune{' ', placeholderRune, '#', ' ', 'a', 'b'},
			expected: []int{-1, -1, 1, -1, -1, -1},
		},
		{
			name:     "h2",
			chars:    []rune{placeholderRune, '#', placeholderRune, '#', ' ', 'a', 'b'},
			expected: []int{-1, 3, -1, -1, -1, -1, -1},
		},
		{
			name:     "h4",
			chars:    []rune{placeholderRune, '#', placeholderRune, '#', placeholderRune, '#', placeholderRune, '#', ' ', 'a'},
			expected: []int{-1, 7, -1, -1, -1, -1, -1, -1, -1, -1},
		},
		{
			name:     "h6",
			chars:    []rune{placeholderRune, '#', placeholderRune, '#', placeholderRune, '#', placeholderRune, '#', placeholderRune, '#', placeholderRune, '#', ' ', 'a'},
			expected: []int{-1, 11, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		},
		{
			name:     "no h7",
			chars:    []rune{placeholderRune, '#', placeholderRune, '#', placeholderRune, '#', placeholderRune, '#', placeholderRune, '#', placeholderRune, '#', placeholderRune, '#', ' ', 'a'},
			expected: []int{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		},
		{
			name:     "tag",
			chars:    []rune{placeholderRune, '#', 'a'},
			expected: []int{-1, -1, -1},
		},
		{
			name:     "also take empty heading",
			chars:    []rune{placeholderRune, '#', ' ', '\n', 'a'},
			expected: []int{-1, 1, -1, -1, -1},
		},
		{
			name:     "nothing afterwards",
			chars:    []rune{placeholderRune, '#'},
			expected: []int{-1, 1},
		},
		{
			name:     "tab before content",
			chars:    []rune{placeholderRune, '#', '\t', 'a'},
			expected: []int{-1, 1, -1, -1},
		},
	}
	for _, run := range runs {

		t.Run(run.name, func(t *testing.T) {

			var actual []int
			for index, _ := range run.chars {
				b := []byte(string(run.chars))
				output := IsAtxHeader(b, index)

				actual = append(actual, output)
			}

			if !reflect.DeepEqual(actual, run.expected) {
				t.Errorf("expected %+v but got %+v", run.expected, actual)
			}
		})

	}
}

func TestIsSetextHeader(t *testing.T) {

	runs := []struct {
		name  string
		chars []rune

		expected []int
	}{
		{
			name:     "not needed",
			chars:    []rune{'a', 'b', 'c'},
			expected: []int{-1, -1, -1},
		},
		{
			name:     "inside text",
			chars:    []rune{'a', placeholderRune, '=', 'b'},
			expected: []int{-1, -1, -1, -1},
		},
		{
			name:     "blank line before",
			chars:    []rune{'a', '\n', '\n', placeholderRune, '='},
			expected: []int{-1, -1, -1, -1, -1},
		},
		{
			name:     "with heading",
			chars:    []rune{'a', '\n', placeholderRune, '='},
			expected: []int{-1, -1, -1, 1},
		},
		{
			name:     "at start of file",
			chars:    []rune{placeholderRune, '='},
			expected: []int{-1, -1},
		},
	}
	for _, run := range runs {
		t.Run(run.name, func(t *testing.T) {

			var actual []int
			for index, _ := range run.chars {
				b := []byte(string(run.chars))
				output := IsSetextHeader(b, index)

				actual = append(actual, output)
			}

			if !reflect.DeepEqual(actual, run.expected) {
				t.Errorf("expected %+v but got %+v", run.expected, actual)
			}
		})
	}
}
