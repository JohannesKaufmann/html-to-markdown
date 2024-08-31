package escape

import (
	"reflect"
	"testing"
)

func TestIsImageOrLink(t *testing.T) {
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
			name:     "image A",
			chars:    []byte{placeholderByte, '!', placeholderByte, '[', 'a', ']'},
			expected: []int{-1, -1, -1, 1, -1, -1},
		},
		{
			name:     "image B",
			chars:    []byte{'!', placeholderByte, '[', 'a', ']'},
			expected: []int{-1, -1, 1, -1, -1},
		},
		{
			name:     "image C",
			chars:    []byte{placeholderByte, '!', '[', 'a', ']'},
			expected: []int{-1, 1, 1, -1, -1},
		},
		{
			name:     "multiple starting brackets",
			chars:    []byte{'[', '[', '[', 'a', ']'},
			expected: []int{1, 1, 1, -1, -1},
		},
		{
			name:     "newline in content",
			chars:    []byte{'[', 'a', '\n', ']'},
			expected: []int{-1, -1, -1, -1},
		},
		{
			name:     "at end of file",
			chars:    []byte{'[', 'a'},
			expected: []int{-1, -1},
		},
		{
			name:     "at end of file",
			chars:    []byte{'!'},
			expected: []int{-1},
		},
	}
	for _, run := range runs {

		t.Run(run.name, func(t *testing.T) {
			var actual []int
			for index := range run.chars {
				output := IsImageOrLink(run.chars, index)

				actual = append(actual, output)
			}

			if !reflect.DeepEqual(actual, run.expected) {
				t.Errorf("expected %+v but got %+v", run.expected, actual)
			}
		})

	}
}
