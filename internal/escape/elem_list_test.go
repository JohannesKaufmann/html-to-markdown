package escape

import (
	"reflect"
	"testing"
)

func TestIsUnorderedList(t *testing.T) {
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
			name:     "dash inside text",
			chars:    []byte{'a', '-', ' ', 'b'},
			expected: []int{-1, -1, -1, -1},
		},
		{
			name:     "dash and directly text",
			chars:    []byte{'-', 'a', 'b'},
			expected: []int{-1, -1, -1},
		},
		{
			name:     "two lists",
			chars:    []byte{placeholderByte, '-', ' ', 'a', '\n', placeholderByte, '-', ' ', 'b'},
			expected: []int{-1, 1, -1, -1, -1, -1, 1, -1, -1},
		},
		{
			name:     "space before list",
			chars:    []byte{' ', '-', ' ', 'a'},
			expected: []int{-1, 1, -1, -1},
		},
	}
	for _, run := range runs {

		t.Run(run.name, func(t *testing.T) {
			var actual []int
			for index := range run.chars {
				output := IsUnorderedList(run.chars, index)

				actual = append(actual, output)
			}

			if !reflect.DeepEqual(actual, run.expected) {
				t.Errorf("expected %+v but got %+v", run.expected, actual)
			}
		})

	}
}

func TestIsOrderedList(t *testing.T) {
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
			name:     "simple list",
			chars:    []byte{'1', '.', ' ', 'a'},
			expected: []int{-1, 1, -1, -1},
		},
		{
			name:     "bigger list",
			chars:    []byte{'1', '2', '3', '.', ' ', 'a'},
			expected: []int{-1, -1, -1, 1, -1, -1},
		},
		{
			name:     "inside text",
			chars:    []byte{'a', '1', '.', ' ', 'a'},
			expected: []int{-1, -1, -1, -1, -1},
		},
		{
			name:     "space after dot missing",
			chars:    []byte{'1', '.', 'a'},
			expected: []int{-1, -1, -1},
		},
		{
			name:     "number before dot missing",
			chars:    []byte{'a', '.', 'b'},
			expected: []int{-1, -1, -1},
		},
		{
			name:     "allow space before dot",
			chars:    []byte{' ', '1', '.', ' ', 'a'},
			expected: []int{-1, -1, 1, -1, -1},
		},
		{
			name:     "two lists",
			chars:    []byte{placeholderByte, '1', '.', ' ', 'a', '\n', placeholderByte, '1', '.', ' ', 'b'},
			expected: []int{-1, -1, 1, -1, -1, -1, -1, -1, 1, -1, -1},
		},
	}
	for _, run := range runs {

		t.Run(run.name, func(t *testing.T) {
			var actual []int
			for index := range run.chars {
				output := IsOrderedList(run.chars, index)

				actual = append(actual, output)
			}

			if !reflect.DeepEqual(actual, run.expected) {
				t.Errorf("expected %+v but got %+v", run.expected, actual)
			}
		})

	}
}
