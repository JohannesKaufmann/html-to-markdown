package escape

import (
	"reflect"
	"testing"
)

func TestIsFencedCode(t *testing.T) {
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
			name:     "only two",
			chars:    []byte{placeholderByte, '`', placeholderByte, '`', 'a'},
			expected: []int{-1, -1, -1, -1, -1},
		},
		{
			name:     "other chars before",
			chars:    []byte{'a', ' ', placeholderByte, '`', placeholderByte, '`', placeholderByte, '`', 'a'},
			expected: []int{-1, -1, -1, -1, -1, -1, -1, -1, -1},
		},
		{
			name:     "just beginning",
			chars:    []byte{placeholderByte, '`', placeholderByte, '`', placeholderByte, '`', 'a'},
			expected: []int{-1, 5, -1, -1, -1, -1, -1},
		},
		{
			name:     "just beginning (with tilde)",
			chars:    []byte{placeholderByte, '~', placeholderByte, '~', placeholderByte, '~', 'a'},
			expected: []int{-1, 5, -1, -1, -1, -1, -1},
		},
		{
			name:     "just beginning (with space before)",
			chars:    []byte{' ', placeholderByte, '`', placeholderByte, '`', placeholderByte, '`', 'a'},
			expected: []int{-1, -1, 5, -1, -1, -1, -1, -1},
		},
		{
			name:     "just beginning (with newline before)",
			chars:    []byte{'\n', placeholderByte, '`', placeholderByte, '`', placeholderByte, '`', 'a'},
			expected: []int{-1, -1, 5, -1, -1, -1, -1, -1},
		},
		{
			name:     "simple",
			chars:    []byte{placeholderByte, '`', placeholderByte, '`', placeholderByte, '`', '\n', 'a', '\n', placeholderByte, '`', placeholderByte, '`', placeholderByte, '`'},
			expected: []int{-1, 5, -1, -1, -1, -1, -1, -1, -1, -1, 5, -1, -1, -1, -1},
		},
		{
			name:     "only one end delimiter",
			chars:    []byte{placeholderByte, '`', placeholderByte, '`', placeholderByte, '`', '\n', 'a', '\n', placeholderByte, '`'},
			expected: []int{-1, 5, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		},
	}
	for _, run := range runs {
		t.Run(run.name, func(t *testing.T) {
			var actual []int
			for index := range run.chars {
				output := IsFencedCode(run.chars, index)

				actual = append(actual, output)
			}

			if !reflect.DeepEqual(actual, run.expected) {
				t.Errorf("expected %+v (l:%d) but got %+v (l:%d)", run.expected, len(run.expected), actual, len(actual))
			}

		})
	}
}

func TestIsInlineCode(t *testing.T) {
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
		// {
		// 	name:     "one delimiter inside text, no end delimiter",
		// 	chars:    []byte{'a', '`', 'b', ' ', '\n', '\n', 'a'},
		// 	expected: []int{-1, -1, -1},
		// },
		{
			name:  "simple",
			chars: []byte{'`', 'a', '`'},
			// expected: []int{1, -1, -1},
			expected: []int{1, -1, 1},
		},
		// {
		// 	name:     "without content",
		// 	chars:    []byte{'`', '`'},
		// 	expected: []int{1, -1},
		// },
		// {
		// 	name:     "code inside normal text",
		// 	chars:    []byte{'a', '`', 'b', '`', 'a'},
		// 	expected: []int{-1, 3, -1, -1, -1},
		// },

		// TODO: also nested: ``some `code` text``

	}
	for _, run := range runs {
		t.Run(run.name, func(t *testing.T) {
			var actual []int
			for index := range run.chars {
				output := IsInlineCode(run.chars, index)

				actual = append(actual, output)
			}

			if !reflect.DeepEqual(actual, run.expected) {
				t.Errorf("expected %+v but got %+v", run.expected, actual)
			}
		})
	}
}
