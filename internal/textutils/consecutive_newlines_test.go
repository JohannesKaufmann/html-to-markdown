package textutils

import (
	"bytes"
	"testing"
)

func TestTrimConsecutiveNewlines(t *testing.T) {
	runs := []struct {
		desc     string
		input    []byte
		expected []byte
	}{
		{
			desc:     "empty",
			input:    []byte(""),
			expected: []byte(""),
		},
		{
			desc:     "not needed",
			input:    []byte("normal text"),
			expected: []byte("normal text"),
		},
		{
			desc:     "also not needed",
			input:    []byte("normal\n\ntext"),
			expected: []byte("normal\n\ntext"),
		},

		{
			desc:     "just two newlines",
			input:    []byte("\n\n"),
			expected: []byte("\n\n"),
		},
		{
			desc:     "just three newlines",
			input:    []byte("\n\n\n"),
			expected: []byte("\n\n"),
		},
		{
			desc:     "just four newlines",
			input:    []byte("\n\n\n\n"),
			expected: []byte("\n\n"),
		},

		{
			desc:     "newlines before",
			input:    []byte("\n\n\ntext"),
			expected: []byte("\n\ntext"),
		},
		{
			desc:     "newlines after",
			input:    []byte("text\n\n\n"),
			expected: []byte("text\n\n"),
		},
		{
			desc:     "newlines before and after",
			input:    []byte("\n\n\ntext\n\n\n"),
			expected: []byte("\n\ntext\n\n"),
		},
		{
			desc:     "newlines between",
			input:    []byte("before\n\n\nafter"),
			expected: []byte("before\n\nafter"),
		},
		{
			desc:     "newlines between multiple times",
			input:    []byte("1\n\n\n2\n\n\n3"),
			expected: []byte("1\n\n2\n\n3"),
		},

		{
			desc:     "not needed the first time",
			input:    []byte("abc\n\nabc\n\n\nabc"),
			expected: []byte("abc\n\nabc\n\nabc"),
		},
		{
			desc:     "not needed the second time",
			input:    []byte("abc\n\n\nabc\n\nabc"),
			expected: []byte("abc\n\nabc\n\nabc"),
		},

		{
			desc:     "with special characters",
			input:    []byte("äöü\n\n\näöü"),
			expected: []byte("äöü\n\näöü"),
		},
		{
			desc:     "space at end",
			input:    []byte("a\n\n\nb "),
			expected: []byte("a\n\nb "),
		},
		{
			desc:     "one newline at end",
			input:    []byte("a\n\n\nb\n"),
			expected: []byte("a\n\nb\n"),
		},
		{
			desc:     "two newlines at end",
			input:    []byte("a\n\n\nb\n\n"),
			expected: []byte("a\n\nb\n\n"),
		},
	}

	for _, run := range runs {
		t.Run(run.desc, func(t *testing.T) {
			output := TrimConsecutiveNewlines(run.input)
			if !bytes.Equal(output, run.expected) {
				t.Errorf("expected %q but got %q", string(run.expected), string(output))
			}
		})
	}
}

func TestTrimConsecutiveNewlines_Allocs(t *testing.T) {
	const N = 1000

	avg := testing.AllocsPerRun(N, func() {
		input := []byte("abc")
		output := TrimConsecutiveNewlines(input)
		_ = output
	})
	if avg != 0 {
		t.Errorf("with no newlines there should be no allocations but got %f", avg)
	}

	avg = testing.AllocsPerRun(N, func() {
		input := []byte("abc\n\nabc")
		output := TrimConsecutiveNewlines(input)
		_ = output
	})
	if avg != 0 {
		t.Errorf("with only two newlines there should be no allocations but got %f", avg)
	}

	avg = testing.AllocsPerRun(N, func() {
		input := []byte("abc\n\n\nabc")
		output := TrimConsecutiveNewlines(input)
		_ = output
	})
	if avg != 1 {
		t.Errorf("with trhee newlines there should be 1 allocation but got %f", avg)
	}
}

const Repeat = 10

func BenchmarkTrimConsecutiveNewlines(b *testing.B) {
	runs := []struct {
		desc  string
		input []byte
	}{
		{
			desc:  "not needed",
			input: bytes.Repeat([]byte("normal\n\ntext"), Repeat),
		},
		{
			desc:  "multiple times",
			input: bytes.Repeat([]byte("1\n\n\n2\n\n\n3"), Repeat),
		},
	}

	for _, run := range runs {
		b.Run(run.desc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				TrimConsecutiveNewlines(run.input)
			}
		})
	}
}
