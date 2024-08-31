package escape

import (
	"bytes"
	"testing"
	"unicode"
)

const ChineseRune = '字'

func TestIsSpace(t *testing.T) {

	runs := []struct {
		name     string
		input    byte
		expected bool
	}{
		{
			name:     "empty",
			input:    0,
			expected: false,
		},
		{
			name:     "normal char a",
			input:    'a',
			expected: false,
		},
		{
			name:     "special character ä",
			input:    'ä',
			expected: false,
		},
		{
			name:     "chinese character #1",
			input:    []byte(string(ChineseRune))[0],
			expected: false,
		},
		{
			name:     "chinese character #2",
			input:    []byte(string(ChineseRune))[1],
			expected: false,
		},
		{
			name:     "chinese character #3",
			input:    []byte(string(ChineseRune))[2],
			expected: false,
		},
		{
			name:     "space",
			input:    ' ',
			expected: true,
		},
		{
			name:     "tab",
			input:    '	',
			expected: true,
		},
	}
	for _, run := range runs {
		t.Run(run.name, func(t *testing.T) {
			t.Run("unicode.IsSpace", func(t *testing.T) {
				output := unicode.IsSpace(rune(run.input))
				if output != run.expected {
					t.Errorf("for %s expected %v but got %v", string(run.input), run.expected, output)
				}
			})
			t.Run("escape.IsSpace", func(t *testing.T) {
				output := IsSpace(run.input)
				if output != run.expected {
					t.Errorf("for %s expected %v but got %v", string(run.input), run.expected, output)
				}
			})
		})
	}
}

func TestRune(t *testing.T) {

	chars := []rune{
		' ',
		'\n',
		'\t',
		rune(6), // Acknowledge character
		rune(7), // Bell character

		'!',
		'"',
		'#',
		'$',
		'%',
		'&',
		'\'',
		'(',
		')',
		'*',
		'+',
		',',
		'-',
		'.',
		'/',
		':',
		';',
		'<',
		'=',
		'>',
		'?',
		'@',
		'[',
		'\\',
		']',
		'^',
		'_',
		'`',
		'{',
		'|',
		'}',
		'~',
	}
	for _, char := range chars {
		t.Run(string(char), func(t *testing.T) {
			length := len([]byte(string(char)))

			if length != 1 {
				t.Errorf("got a length of %d", length)
			}

		})
	}
}

func TestGetPrev(t *testing.T) {
	input := []byte{'a', placeholderByte, 'b', 'c'}

	if getPrev(input, 3) != 'b' {
		t.Error("expected different output")
	}
	if getPrev(input, 2) != 'a' {
		t.Error("expected different output")
	}
	if getPrev(input, 1) != 'a' {
		t.Error("expected different output")
	}
	if getPrev(input, 0) != 0 {
		t.Error("expected different output")
	}
}

func TestGetNext(t *testing.T) {
	input := []byte{'a', placeholderByte, 'b', 'c'}

	if getNext(input, 0) != 'b' {
		t.Error("expected different output")
	}
	if getNext(input, 1) != 'b' {
		t.Error("expected different output")
	}
	if getNext(input, 2) != 'c' {
		t.Error("expected different output")
	}
	if getNext(input, 3) != 0 {
		t.Error("expected different output")
	}
}

func TestGetNextAsRune(t *testing.T) {
	inputString := "a\a⌘\ab"
	inputBytes := []byte{
		97,            // a
		7,             // bell (our escape char)
		226, 140, 152, // mac sign
		7,  // bell (our escape char)
		98, // b
	}

	if !bytes.Equal([]byte(inputString), inputBytes) {
		t.Error("the string and byte slice dont match")
	}

	nextByte := getNext(inputBytes, 0)
	if nextByte != 226 {
		t.Error("expected different next byte")
	}

	nextRune := getNextAsRune(inputBytes, 0)
	if nextRune != '⌘' {
		t.Error("expected different next rune")
	}

	lastNextRune := getNextAsRune(inputBytes, 6)
	if lastNextRune != 0 {
		t.Error("expected different last next rune")
	}
	// - - - - //

	prevByte := getPrev(inputBytes, 6)
	if prevByte != 152 {
		t.Error("expected different prev byte")
	}

	prevRune := getPrevAsRune(inputBytes, 6)
	if prevRune != '⌘' {
		t.Error("expected different prev rune")
	}
	firstPrevRune := getPrevAsRune(inputBytes, 0)
	if firstPrevRune != 0 {
		t.Error("expected different zero prev rune")
	}
}
