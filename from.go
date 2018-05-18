package md

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Timeout for the http client
var Timeout = time.Second * 10
var netClient = &http.Client{
	Timeout: Timeout,
}

// DomainFromURL removes the path from the url.
func DomainFromURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	u.Path = ""
	return u.String()
}

// FromSelection returns the content from a goquery selection.
// If you have a goquery document just pass in doc.Selection.
func FromSelection(domain string, selec *goquery.Selection) (string, []*Element) {
	elements := SelecToElem(domain, false, selec, nil)

	data, _ := json.Marshal(elements)
	fmt.Println(string(data))

	return ElemToMD(Root, elements), elements
}

// FromString returns the content from a html string. If you
// already have a goquery selection use `FromSelection`.
func FromString(domain, html string) (string, []*Element, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", nil, nil
	}
	markdown, elements := FromSelection(domain, doc.Selection)
	return markdown, elements, nil
}

// FromURL returns the content from the page with that url.
func FromURL(url string) (string, []*Element, error) {
	// not using goquery.NewDocument directly because of the timeout
	resp, err := netClient.Get(url)
	if err != nil {
		return "", nil, nil
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return "", nil, nil
	}
	markdown, elements := FromSelection(DomainFromURL(url), doc.Selection)
	return markdown, elements, nil
}
