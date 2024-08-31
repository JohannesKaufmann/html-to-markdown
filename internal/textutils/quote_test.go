package textutils

import (
	"reflect"
	"testing"
)

func TestSurroundByQuotes(t *testing.T) {

	tests := []struct {
		name  string
		input []byte
		want  []byte
	}{
		{
			name:  "empty content",
			input: nil,
			want:  nil,
		},
		{
			name:  "no content",
			input: []byte(""),
			want:  nil,
		},
		{
			name:  "contains no quotes",
			input: []byte(`no quotes are here`),
			want:  []byte(`"no quotes are here"`),
		},
		{
			name:  "contains double quotes",
			input: []byte(`double "quotes" are here`),
			want:  []byte(`'double "quotes" are here'`),
		},
		{
			name:  "contains single quotes",
			input: []byte(`single 'quotes' are here`),
			want:  []byte(`"single 'quotes' are here"`),
		},
		{
			name:  "contains both quotes",
			input: []byte(`double " AND single ' quotes '" are here`),
			want:  []byte(`"double \" AND single ' quotes '\" are here"`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SurroundByQuotes(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SurroundByQuotes() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
