package textutils

import "testing"

func TestDelimiterForEveryLine(t *testing.T) {
	tests := []struct {
		name string

		text      string
		delimiter string

		want string
	}{
		{
			name: "put delimiter around text",

			text:      "bold text",
			delimiter: "**",

			want: "**bold text**",
		},
		{
			name: "keep whitespace outside (normal space)",

			text:      " bold text ",
			delimiter: "**",

			want: " **bold text** ",
		},
		{
			name: "keep whitespace outside (non-breaking space)",

			text:      "\u00a0bold text\u00a0\u00a0",
			delimiter: "**",

			want: "\u00a0**bold text**\u00a0\u00a0",
		},
		{
			name: "keep whitespace outside on every line (non-breaking space)",

			text:      "bold\u00a0\ntext\u00a0",
			delimiter: "**",

			want: "**bold**\u00a0\n**text**\u00a0",
		},
		{
			name: "put strong on every line",

			text:      "line 1\nline 2",
			delimiter: "**",

			want: "**line 1**\n**line 2**",
		},
		{
			name: "skip empty lines",

			text:      "line 1\n\n\nline 2",
			delimiter: "_",

			want: "_line 1_\n\n\n_line 2_",
		},
		{
			name: "with indentation",

			text: `
line 1

line 2
line 3
`,
			delimiter: "__",

			want: `
__line 1__

__line 2__
__line 3__
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := DelimiterForEveryLine([]byte(tt.text), []byte(tt.delimiter)); string(got) != tt.want {
				t.Errorf("DelimiterForEveryLine() = \n'%v' but want \n'%v'", string(got), tt.want)
			}

		})
	}
}
