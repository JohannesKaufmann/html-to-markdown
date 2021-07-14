package md

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
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

type ErrReader struct{ Error error }

// -> https://stackoverflow.com/a/57452918
func (e *ErrReader) Read([]byte) (int, error) {
	return 0, e.Error
}
func TestConvertReader_Error(t *testing.T) {
	reader := &ErrReader{
		Error: errors.New("we got an error"),
	}

	converter := NewConverter("", true, nil)
	res, err := converter.ConvertReader(reader)
	if err != reader.Error {
		t.Error("expected an error")
	}

	if res.Len() != 0 {
		t.Error("expected an empty buffer")
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

func TestConvertBytes_Empty(t *testing.T) {
	converter := NewConverter("", true, nil)
	res, err := converter.ConvertBytes(nil)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(res, []byte("")) {
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

func TestConvertResponse_Error(t *testing.T) {
	expectedErr := errors.New("custom error reader")

	converter := NewConverter("", true, nil)
	res, err := converter.ConvertResponse(&http.Response{
		StatusCode: 200,
		Body: ioutil.NopCloser(&ErrReader{
			Error: expectedErr,
		}),
		Request: &http.Request{},
	})
	if err != expectedErr {
		t.Error(err)
	}

	if res != "" {
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

func TestConvertURL_Error(t *testing.T) {
	url := "abc https://example.com"

	converter := NewConverter("", true, nil)
	res, err := converter.ConvertURL(url)
	if err == nil {
		t.Error("expected an error")
	}

	if res != "" {
		t.Error("the result is different that expected")
	}
}

func TestConvertURL_ErrorStatusCode(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("404 Not Found"))
	}))
	// Close the server when test finishes
	defer server.Close()
	// override the client used in `ConvertURL`
	netClient = server.Client()

	converter := NewConverter(server.URL, true, nil)
	res, err := converter.ConvertURL(server.URL)
	if err == nil {
		t.Error("expected an error")
	}

	if res != "" {
		t.Error("the result is different that expected")
	}
}

// - - - - - - - - - - - - //

func TestNewConverter_NoRules(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetFlags(0)
	defer func() {
		// reset the options back to the defaults
		log.SetOutput(os.Stderr)
		log.SetFlags(3)
	}()

	input := `<strong>Bold</strong>`
	expected := ``

	// disable commonmark
	converter := NewConverter("", false, nil)
	res, err := converter.ConvertString(input)
	if err != nil {
		t.Error(err)
	}

	if res != expected {
		t.Error("the result is different that expected")
	}

	if strings.TrimSuffix(buf.String(), "\n") != "you have added no rules. either enable commonmark or add you own." {
		t.Error("expected a different log message")
	}
}

func TestNewConverter_ValidateOptions(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetFlags(0)
	defer func() {
		// reset the options back to the defaults
		log.SetOutput(os.Stderr)
		log.SetFlags(3)
	}()

	input := `<strong>Bold</strong>`
	expected := `====Bold====`

	converter := NewConverter("", true, &Options{
		StrongDelimiter: "====",
	})
	res, err := converter.ConvertString(input)
	if err != nil {
		t.Error(err)
	}

	if res != expected {
		t.Error("the result is different that expected")
	}

	if strings.TrimSuffix(buf.String(), "\n") != "markdown options is not valid: field must be one of [** __] but got ====" {
		t.Error("expected a different log message")
	}
}

func TestNewConverter_ValidateOptions_All(t *testing.T) {
	var tests = []struct {
		name    string
		options *Options

		input    string
		expected string
	}{
		{
			name: "HeadingStyle",
			options: &Options{
				HeadingStyle: "invalid",
			},
			input:    `<h1>Heading</h1>`,
			expected: `# Heading`,
		},
		{
			name: "HorizontalRule",
			options: &Options{
				HorizontalRule: "--",
			},
			input:    `<hr />`,
			expected: `--`,
		},
		{
			name: "BulletListMarker",
			options: &Options{
				BulletListMarker: "^",
			},
			input:    `<ul><li>Test</li></ul>`,
			expected: `^ Test`,
		},
		{
			name: "CodeBlockStyle",
			options: &Options{
				CodeBlockStyle: "invalid",
			},
			input:    `<code>test</code>`,
			expected: "`test`",
		},
		{
			name: "Fence",
			options: &Options{
				Fence: "^^^",
			},
			input:    `<pre>test</pre>`,
			expected: "^^^\ntest\n^^^",
		},
		{
			name: "EmDelimiter",
			options: &Options{
				EmDelimiter: "-",
			},
			input:    `<i>test</i>`,
			expected: "-test-",
		},
		{
			name: "LinkStyle",
			options: &Options{
				LinkStyle: "invalid",
			},
			input: `<a href="example.com">link</a>`,
			expected: `[link][1]

[1]: example.com`,
		},
		{
			name: "LinkReferenceStyle",
			options: &Options{
				LinkReferenceStyle: "invalid",
			},
			input:    `<a href="example.com">link</a>`,
			expected: "[link](example.com)",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			log.SetFlags(0)
			defer func() {
				// reset the options back to the defaults
				log.SetOutput(os.Stderr)
				log.SetFlags(3)
			}()

			converter := NewConverter("", true, test.options)
			res, err := converter.ConvertString(test.input)
			if err != nil {
				t.Error(err)
			}

			if res != test.expected {
				t.Errorf("expected '%s' but got '%s'", test.expected, res)
			}

			logOutput := strings.TrimSuffix(buf.String(), "\n")
			if !strings.Contains(logOutput, "markdown options is not valid: ") {
				t.Errorf("expected a different log message but got '%s'", logOutput)
			}
		})
	}

}

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

func TestAddRules_NoRules(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetFlags(0)
	defer func() {
		// reset the options back to the defaults
		log.SetOutput(os.Stderr)
		log.SetFlags(3)
	}()

	var wasCalled bool
	rule := Rule{
		Filter: []string{ /* nothing */ },
		Replacement: func(content string, selec *goquery.Selection, opt *Options) *string {
			wasCalled = true
			return nil
		},
	}

	conv := NewConverter("", true, nil)
	conv.AddRules(rule)

	md, err := conv.ConvertString(`<p>Some Content</p>`)
	if err != nil {
		t.Error(err)
	}
	if md != "Some Content" {
		t.Error("got different markdown result")
	}

	logOutput := strings.TrimSuffix(buf.String(), "\n")
	if logOutput != "you need to specify at least one filter for your rule" {
		t.Errorf("expected a different log message but got '%s'", logOutput)
	}

	if wasCalled {
		t.Error("the rule should not have been called")
	}
}

func TestBefore(t *testing.T) {
	var firstWasCalled bool
	var secondWasCalled bool
	firstHook := func(selec *goquery.Selection) {
		firstWasCalled = true
		if secondWasCalled {
			t.Error("the second hook should not be called yet")
		}
	}
	secondHook := func(selec *goquery.Selection) {
		secondWasCalled = true
		if !firstWasCalled {
			t.Error("the first hook should already be called")
		}
	}

	conv := NewConverter("", true, nil)
	conv.Before(firstHook, secondHook)
	_, err := conv.ConvertString(`<a href="https://test.com">Link</a>`)
	if err != nil {
		t.Error(err)
	}

	if !firstWasCalled || !secondWasCalled {
		t.Error("not all hooks were called")
	}
}

func TestAfter(t *testing.T) {
	var firstWasCalled bool
	var secondWasCalled bool
	firstHook := func(md string) string {
		firstWasCalled = true
		if secondWasCalled {
			t.Error("the second hook should not be called yet")
		}
		return md + " first"
	}
	secondHook := func(md string) string {
		secondWasCalled = true
		if !firstWasCalled {
			t.Error("the first hook should already be called")
		}
		return md + " second"
	}

	conv := NewConverter("", true, nil)
	conv.After(firstHook, secondHook)
	md, err := conv.ConvertString(`<span>base</span>`)
	if err != nil {
		t.Error(err)
	}

	if md != `base first second` {
		t.Errorf("expected different markdown result but got '%s'", md)
	}

	if !firstWasCalled || !secondWasCalled {
		t.Error("not all hooks were called")
	}
}
func TestClearBefore(t *testing.T) {
	var wasCalled bool
	hook := func(selec *goquery.Selection) {
		wasCalled = true
	}

	conv := NewConverter("", true, nil)
	conv.ClearBefore()
	if len(conv.before) != 0 {
		t.Error("the before hook array should be of length 0")
	}

	conv.Before(hook)

	_, err := conv.ConvertString(`<a href="https://test.com">Link</a>`)
	if err != nil {
		t.Error(err)
	}

	if !wasCalled {
		t.Error("the hook should have been called")
	}
}

func TestClearAfter(t *testing.T) {
	var wasCalled bool
	hook := func(markdown string) string {
		wasCalled = true
		return "my new value"
	}

	conv := NewConverter("", true, nil)
	conv.ClearAfter()
	if len(conv.after) != 0 {
		t.Error("the after hook array should be of length 0")
	}

	conv.After(hook)

	md, err := conv.ConvertString(`<a href="https://test.com">Link</a>`)
	if err != nil {
		t.Error(err)
	}
	if md != "my new value" {
		t.Error("the result was different then expected")
	}

	if !wasCalled {
		t.Error("the hook should have been called")
	}
}

func TestDomainFromURL(t *testing.T) {
	var tests = []struct {
		input    string
		expected string
	}{
		{
			input:    "example.com",
			expected: "example.com",
		},
		{
			input:    "https://example.com",
			expected: "example.com",
		},
		{
			input:    "https://www.example.com",
			expected: "www.example.com",
		},

		{
			input:    "http://example.com/index.html",
			expected: "example.com",
		},
		{
			input:    "http://example.com?page=home",
			expected: "example.com",
		},
		{
			input:    "http://example.com#page",
			expected: "example.com",
		},
		{
			input:    "http://example.com:3000",
			expected: "example.com:3000",
		},
		{
			// not so happy about this :(
			input:    "example",
			expected: "example",
		},
		{
			input:    "https://developer.mozilla.org/en-US/docs/Web/API/URL/host",
			expected: "developer.mozilla.org",
		},
		{
			input:    "  http://example.com",
			expected: "example.com",
		},
		{
			// invalid url
			input:    "abc  http://example.com",
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			res := DomainFromURL(test.input)
			if res != test.expected {
				t.Errorf("for '%s' expected '%s' but got '%s'", test.input, test.expected, res)
			}
		})
	}
}
