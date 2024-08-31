package cmd

import (
	"reflect"
	"testing"
)

func TestFlagStringSlice(t *testing.T) {
	testCases := []struct {
		desc     string
		inputs   []string
		expected []string
	}{
		{
			desc:     "simple flag",
			inputs:   []string{"a"},
			expected: []string{"a"},
		},
		{
			desc:     "two flags",
			inputs:   []string{"a,b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			desc:     "with seperator",
			inputs:   []string{"a,", ",b"},
			expected: []string{"a", "b"},
		},
		{
			desc:     "with spaces",
			inputs:   []string{"a, ,b", " ,c"},
			expected: []string{"a", "b", "c"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			var result []string
			for _, input := range tC.inputs {
				flagStringSlice(&result)(input)
			}

			if !reflect.DeepEqual(result, tC.expected) {
				t.Errorf("expected %v but got %v", tC.expected, result)
			}
		})
	}
}
