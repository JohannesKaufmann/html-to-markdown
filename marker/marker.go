package marker

import (
	"bytes"
	"unicode"
	"unicode/utf8"
)

const (
	// For simplificity we are using a rune that is one byte wide. A character
	// that is not used widely (apart from cli's) is the bell character (7).
	MarkerEscaping rune = '\a'

	// - - - - //

	// Marker1                rune = '\uF000' // 61440
	MarkerLineBreak        rune = '\uF001' // 61441
	MarkerCodeBlockNewline rune = '\uF002' // 61442
)

var (
	BytesMarkerEscaping = []byte{7}

	// BytesMarker1                = []byte{239, 128, 128}
	BytesMarkerLineBreak        = []byte{239, 128, 129}
	BytesTWICEMarkerLineBreak   = []byte{239, 128, 129, 239, 128, 129}
	BytesMarkerCodeBlockNewline = []byte{239, 128, 130}
)

func init() {
	checkRuneAndByteSlice(MarkerEscaping, BytesMarkerEscaping)
	checkRuneAndByteSlice(MarkerLineBreak, BytesMarkerLineBreak)
	checkRuneAndByteSlice(MarkerCodeBlockNewline, BytesMarkerCodeBlockNewline)
}

func checkRuneAndByteSlice(r rune, b []byte) {
	if !bytes.Equal([]byte(string(r)), b) {
		panic("the rune and byte slice dont represent the same character")
	}
}

func GetMarker(p []byte, i int) (marker rune, size int) {
	r, size := utf8.DecodeRune(p[i:])

	switch r {
	case MarkerLineBreak, MarkerCodeBlockNewline:
		return r, size

	default:
		return 0, 0
	}
}

func IsSpace(r rune) bool {
	return unicode.IsSpace(r) || r == MarkerLineBreak
}

// func IsNewline(r rune) bool {
// 	return r == '\n' || r == '\r' || r == MarkerLineBreak
// }

// TODO: should this be in another package?
func SplitFunc(str []byte, fn func(rune) bool) [][]byte {
	var substrs [][]byte
	for {
		i := bytes.IndexFunc(str, fn)
		if i == -1 {
			if len(str) > 0 {
				substrs = append(substrs, str)
			}
			break
		}

		_, size := utf8.DecodeRune(str[i:])
		// substrs = append(substrs, str[:i], str[i:i+1])
		substrs = append(substrs, str[:i])
		str = str[i+size:]
	}

	return substrs
}
