package md

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestConvertReader(t *testing.T) {
	input := `<strong>Bold</strong>`
	expected := `**Bold**`

	converter := NewConverter("", true, nil)
	res, err := converter.ConvertReader(strings.NewReader(input))
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(res.Bytes(), []byte(expected)) {
		t.Error("the result is different that expected")
	}
}

func TestConvertBytes(t *testing.T) {
	input := `<strong>Bold</strong>`
	expected := `**Bold**`

	converter := NewConverter("", true, nil)
	res, err := converter.ConvertBytes([]byte(input))
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(res, []byte(expected)) {
		t.Error("the result is different that expected")
	}
}

func TestConvertResponse(t *testing.T) {
	input := `<strong>Bold</strong>`
	expected := `**Bold**`

	converter := NewConverter("", true, nil)
	res, err := converter.ConvertResponse(&http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString(input)),
		Request:    &http.Request{},
	})
	if err != nil {
		t.Error(err)
	}

	if res != expected {
		t.Error("the result is different that expected")
	}
}

func TestConvertString(t *testing.T) {
	input := `<strong>Bold</strong>`
	expected := `**Bold**`

	converter := NewConverter("", true, nil)
	res, err := converter.ConvertString(input)
	if err != nil {
		t.Error(err)
	}

	if res != expected {
		t.Error("the result is different that expected")
	}
}

func TestConvertSelection(t *testing.T) {
	input := `<strong>Bold</strong>`
	expected := `**Bold**`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(input))
	if err != nil {
		t.Error(err)
	}

	converter := NewConverter("", true, nil)
	res := converter.Convert(doc.Selection)

	if res != expected {
		t.Error("the result is different that expected")
	}
}

func TestConvertURL(t *testing.T) {
	input := `<strong>Bold</strong>`
	expected := `**Bold**`

	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(input))
	}))
	// Close the server when test finishes
	defer server.Close()
	// override the client used in `ConvertURL`
	netClient = server.Client()

	converter := NewConverter(server.URL, true, nil)
	res, err := converter.ConvertURL(server.URL)
	if err != nil {
		t.Error(err)
	}

	if res != expected {
		t.Error("the result is different that expected")
	}
}

// - - - - - - - - - - - - //

func BenchmarkFromString(b *testing.B) {
	converter := NewConverter("www.google.com", true, nil)

	strongRule := Rule{
		Filter: []string{"strong"},
		Replacement: func(content string, selec *goquery.Selection, opt *Options) *string {
			return nil
		},
	}

	var wg sync.WaitGroup
	convert := func(html string) {
		defer wg.Done()
		_, err := converter.ConvertString(html)
		if err != nil {
			b.Error(err)
		}
	}
	add := func() {
		defer wg.Done()
		converter.AddRules(strongRule)
	}

	for n := 0; n < b.N; n++ {
		wg.Add(2)
		go add()
		go convert("<strong>Bold</strong>")
	}

	wg.Wait()
}

func TestAddRules_ChangeContent(t *testing.T) {
	expected := "Some other Content"

	var wasCalled bool
	rule := Rule{
		Filter: []string{"p"},
		Replacement: func(content string, selec *goquery.Selection, opt *Options) *string {
			wasCalled = true

			if content != "Some Content" {
				t.Errorf("got wrong `content`: '%s'", content)
			}
			if !selec.Is("p") {
				t.Error("selec is not p")
			}
			return String(expected)
		},
	}

	conv := NewConverter("", true, nil)
	conv.AddRules(rule)
	md, err := conv.ConvertString(`<p>Some Content</p>`)
	if err != nil {
		t.Error(err)
	}

	if md != expected {
		t.Errorf("wanted '%s' but got '%s'", expected, md)
	}
	if !wasCalled {
		t.Error("rule was not called")
	}
}

func TestAddRules_Fallback(t *testing.T) {
	// firstExpected := "Some other Content"
	expected := "Totally different Content"

	var firstWasCalled bool
	var secondWasCalled bool
	firstRule := Rule{
		Filter: []string{"p"},
		Replacement: func(content string, selec *goquery.Selection, opt *Options) *string {
			firstWasCalled = true
			if secondWasCalled {
				t.Error("expected first rule to be called before second rule. second is already called")
			}

			if content != "Some Content" {
				t.Errorf("got wrong `content`: '%s'", content)
			}
			if !selec.Is("p") {
				t.Error("selec is not p")
			}

			return nil
		},
	}
	secondRule := Rule{
		Filter: []string{"p"},
		Replacement: func(content string, selec *goquery.Selection, opt *Options) *string {
			secondWasCalled = true
			if !firstWasCalled {
				t.Error("expected first rule to be called before second rule. first is not called yet")
			}

			if content != "Some Content" {
				t.Errorf("got wrong `content`: '%s'", content)
			}
			if !selec.Is("p") {
				t.Error("selec is not p")
			}
			return String(expected)
		},
	}

	conv := NewConverter("", true, nil)
	conv.AddRules(secondRule, firstRule)
	md, err := conv.ConvertString(`<p>Some Content</p>`)
	if err != nil {
		t.Error(err)
	}

	if md != expected {
		t.Errorf("wanted '%s' but got '%s'", expected, md)
	}
	if !firstWasCalled {
		t.Error("first rule was not called")
	}
	if !secondWasCalled {
		t.Error("second rule was not called")
	}
}
