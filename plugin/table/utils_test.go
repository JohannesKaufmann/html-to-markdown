package table

import (
	"reflect"
	"testing"
)

func TestCalculateMaxCounts(t *testing.T) {
	a := [][][]byte{
		{
			[]byte("Company A"),  //  9
			[]byte("Max Müller"), // 10 <--
			[]byte("Berlin"),     //  6 <--
		},
		{
			[]byte("Company Example"), // 15 <--
			[]byte("John Doe"),        //  8
			[]byte("Bonn"),            //  4
		},
		{
			[]byte("A"),
		},
	}

	output := calculateMaxCounts(a)
	expected := []int{15, 10, 6}

	if !reflect.DeepEqual(output, expected) {
		t.Errorf("expected %+v but got %v", expected, output)
	}
}
func TestFillUpRows(t *testing.T) {
	input := [][][]byte{
		{
			[]byte("Company A"),
			[]byte("Max Müller"),
			[]byte("Berlin"),
		},
		{
			[]byte("Company Example"),
			[]byte("John Doe"),
			[]byte("Bonn"),
		},
		{
			[]byte("A"),
			// <--
			// <--
		},
	}

	counts := calculateMaxCounts(input)
	t.Log("counts:", counts)

	// - - - - - - - - - - - - - - - - - - - - //
	maxColumnCount := len(counts)

	output := fillUpRows(input, maxColumnCount)
	expected := [][][]byte{
		{
			[]byte("Company A"),
			[]byte("Max Müller"),
			[]byte("Berlin"),
		},
		{
			[]byte("Company Example"),
			[]byte("John Doe"),
			[]byte("Bonn"),
		},
		{
			[]byte("A"),
			[]byte(""),
			[]byte(""),
		},
	}

	if !reflect.DeepEqual(output, expected) {
		t.Errorf("expected %+v but got %v", expected, output)
	}
}

func TestCalculateModifications(t *testing.T) {
	testCases := []struct {
		desc string

		currentRowIndex int
		currentColIndex int
		colSpan         int
		rowSpan         int

		expected []modification
	}{
		{
			desc: "no modifications needed #1",

			currentRowIndex: 0,
			currentColIndex: 0,
			colSpan:         1,
			rowSpan:         1,

			expected: []modification{},
		},
		{
			desc: "no modifications needed #2",

			currentRowIndex: 10,
			currentColIndex: 5,
			colSpan:         1,
			rowSpan:         1,

			expected: []modification{},
		},

		{
			desc: "colspan=2",

			currentRowIndex: 0,
			currentColIndex: 0,
			colSpan:         2,
			rowSpan:         1,

			expected: []modification{{y: 0, x: 1}},
		},
		{
			desc: "rowspan=2",

			currentRowIndex: 0,
			currentColIndex: 0,
			colSpan:         1,
			rowSpan:         2,

			expected: []modification{{y: 1, x: 0}},
		},
		{
			desc: "colspan=2 and rowspan=2",

			currentRowIndex: 0,
			currentColIndex: 0,
			colSpan:         2,
			rowSpan:         2,

			expected: []modification{
				/* the actual cell  */ {y: 0, x: 1},
				{y: 1, x: 0}, {y: 1, x: 1},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual := calculateModifications(tC.currentRowIndex, tC.currentColIndex, tC.rowSpan, tC.colSpan)
			if len(actual) != len(tC.expected) {
				t.Errorf("expected length %d but got %d", len(tC.expected), len(actual))
			}

			if !reflect.DeepEqual(actual, tC.expected) {
				t.Errorf("expected %+v but got %+v", tC.expected, actual)
			}
		})
	}
}

func TestApplyModifications(t *testing.T) {
	testCases := []struct {
		desc string

		contents      [][][]byte
		modifications []modification

		expected [][][]byte
	}{
		{
			desc: "add in same row",

			contents: [][][]byte{
				{
					[]byte("A"),
				},
			},
			modifications: []modification{
				{
					y: 0,
					x: 0,
				},
			},

			expected: [][][]byte{
				{
					[]byte(""),
					[]byte("A"),
				},
			},
		},
		{
			desc: "add in row below",

			contents: [][][]byte{
				{
					[]byte("A"),
				},
				{
					[]byte("B"),
				},
			},
			modifications: []modification{
				{
					y: 1,
					x: 0,
				},
			},

			expected: [][][]byte{
				{
					[]byte("A"),
				},
				{
					[]byte(""),
					[]byte("B"),
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			applyModifications(tC.contents, tC.modifications)

			if !reflect.DeepEqual(tC.contents, tC.expected) {
				t.Errorf("expected %+v but got %+v", tC.expected, tC.contents)
			}
		})
	}
}
