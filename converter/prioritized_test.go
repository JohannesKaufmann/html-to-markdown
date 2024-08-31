package converter

import (
	"reflect"
	"testing"
)

func TestPrioritizedSlice(t *testing.T) {

	var values = prioritizedSlice[string]{
		prioritized("b", PriorityStandard),
		prioritized("c", PriorityLate),
		prioritized("a", PriorityEarly),
	}

	values.Sort()

	var expected = prioritizedSlice[string]{
		{
			Value:    "a",
			Priority: PriorityEarly,
		},
		{
			Value:    "b",
			Priority: PriorityStandard,
		},
		{
			Value:    "c",
			Priority: PriorityLate,
		},
	}

	if !reflect.DeepEqual(values, expected) {
		t.Errorf("expected %+v but got %+v", expected, values)
	}
}
