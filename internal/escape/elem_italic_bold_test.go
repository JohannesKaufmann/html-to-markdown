package escape

import (
	"reflect"
	"testing"
)

func TestIsItalicOrBold(t *testing.T) {
	runs := []struct {
		desc  string
		chars []rune
		index int

		expected int
	}{
		{
			desc:  "not needed",
			chars: []rune("normal text"),
			index: 0,

			expected: -1,
		},

		{
			desc:  "nothing before",
			chars: []rune("*a"),
			index: 0,

			expected: 1,
		},
		{
			desc:  "newline before",
			chars: []rune("\n *a"),
			index: 2,

			expected: 1,
		},
		{
			desc:  "text and space before",
			chars: []rune("text *a"),
			index: 5,

			expected: 1,
		},

		{
			desc:  "character directly before",
			chars: []rune("a*a"),
			index: 1,

			expected: 1,
		},
		{
			desc:  "point before",
			chars: []rune(".*a"),
			index: 1,

			expected: 1,
		},

		// - - - - //
		{
			desc:     "char after",
			chars:    []rune(" *a"),
			index:    1,
			expected: 1,
		},
		{
			desc:     "point after",
			chars:    []rune(" *."),
			index:    1,
			expected: 1,
		},
		{
			desc:     "nothing after",
			chars:    []rune(" *"),
			index:    1,
			expected: -1,
		},
		{
			desc:     "space after",
			chars:    []rune(" * "),
			index:    1,
			expected: -1,
		},
		{
			desc:     "newline after",
			chars:    []rune(" *\n"),
			index:    1,
			expected: -1,
		},

		// - - - - //
		// {
		// 	desc:     "char before & point after",
		// 	chars:    []rune("a*."),
		// 	index:    1,
		// 	expected: -1,
		// },
		{
			desc:     "space before & point after",
			chars:    []rune(" *."),
			index:    1,
			expected: 1,
		},
		{
			desc:     "point before & point after",
			chars:    []rune(".*."),
			index:    1,
			expected: 1,
		},

		// - - - - //
		{
			desc:     "exclamation mark as content",
			chars:    []rune("*!*"),
			index:    0,
			expected: 1,
		},
		{
			desc:     "special char before",
			chars:    []rune("$*content*"),
			index:    1,
			expected: 1,
		},
		{
			desc:     "$ before",
			chars:    []rune("0$*!*"),
			index:    2,
			expected: 1,
		},
		// {
		// 	desc:     "\x00 after",
		// 	chars:    []rune("*\x00*"),
		// 	index:    0,
		// 	expected: 1,
		// },
		// {
		// 	desc:     "some random input #1",
		// 	chars:    []rune("\xac\xac\xac*!0*"),
		// 	index:    3,
		// 	expected: 1,
		// },
		// {
		// 	desc:     "random input #1",
		// 	chars:    []rune{'\xac', '*', 'a', '*'}, //"*!0*"),
		// 	index:    2,
		// 	expected: 1,
		// },
	}

	for _, run := range runs {
		t.Run(run.desc, func(t *testing.T) {
			b := []byte(string(run.chars))

			t.Log("input:", string(b))
			match := IsItalicOrBold(b, run.index)
			if match != run.expected {
				t.Errorf("expected %d but got %d", run.expected, match)
			}
		})
	}

}

func TestIsItalicOrBold_All(t *testing.T) {
	runs := []struct {
		name  string
		chars []rune

		expected []int
	}{
		{
			name:     "random input #1",
			chars:    []rune{'\xac', '*', '!', '0', '*'},
			expected: []int{-1, -1, 1, -1, -1, -1},
		},
		{
			name:     "random input #2",
			chars:    []rune{'\xac', '\xac', '\xac', '*', '!', '0', '*'},
			expected: []int{-1, -1, -1, -1, -1, -1, 1, -1, -1, -1},
		},
	}
	for _, run := range runs {

		t.Run(run.name, func(t *testing.T) {

			bytes := []byte(string(run.chars))

			var actual []int
			for index := range bytes {
				output := IsItalicOrBold(bytes, index)

				actual = append(actual, output)
			}

			if !reflect.DeepEqual(actual, run.expected) {
				t.Errorf("expected %+v but got %+v", run.expected, actual)
			}
		})

	}

}
