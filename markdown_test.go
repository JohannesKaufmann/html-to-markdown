package md

import "testing"

func TestDefaultGetAbsoluteURL_NoDomain(t *testing.T) {
	input := "/page.html?key=val#hash"
	expected := input

	res := DefaultGetAbsoluteURL(nil, input, "")
	if res != expected {
		t.Errorf("expected '%s' but got '%s'", expected, res)
	}
}

func TestDefaultGetAbsoluteURL_Domain(t *testing.T) {
	input := "/page.html?key=val#hash"
	expected := "http://test.com" + input

	res := DefaultGetAbsoluteURL(nil, input, "test.com")
	if res != expected {
		t.Errorf("expected '%s' but got '%s'", expected, res)
	}
}

func TestDefaultGetAbsoluteURL_DataURI(t *testing.T) {
	input := "data:image/gif;base64,R0lGODlhEAAQAMQAAORHHOVSKudfOulrSOp3WOyDZu6QdvCchPGolfO0o/XBs/fNwfjZ0frl3/zy7////wAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACH5BAkAABAALAAAAAAQABAAAAVVICSOZGlCQAosJ6mu7fiyZeKqNKToQGDsM8hBADgUXoGAiqhSvp5QAnQKGIgUhwFUYLCVDFCrKUE1lBavAViFIDlTImbKC5Gm2hB0SlBCBMQiB0UjIQA7"
	expected := input

	res := DefaultGetAbsoluteURL(nil, input, "test.com")
	if res != expected {
		t.Errorf("expected '%s' but got '%s'", expected, res)
	}
}
