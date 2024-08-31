package marker

import (
	"bytes"
	"reflect"
	"testing"
)

func TestSplitFunc_Basics(t *testing.T) {
	in := []byte("one-two-three")

	out1 := bytes.Split(in, []byte("-"))
	out2 := SplitFunc(in, func(r rune) bool {
		return r == '-'
	})

	if !reflect.DeepEqual(out1, out2) {
		t.Error("the split functions generated different outputs")
	}
}
func TestSplitFunc_SpecialChars(t *testing.T) {
	in := []byte("[öü]ä[öü]ä[öü]")

	out1 := bytes.Split(in, []byte("ä"))
	out2 := SplitFunc(in, func(r rune) bool {
		return r == 'ä'
	})

	if !reflect.DeepEqual(out1, out2) {
		t.Error("the split functions generated different outputs")
	}
}
