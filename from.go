package md

import (
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

// var newLinesR = regexp.MustCompile(``)

// FromSelection returns the content from a goquery selection.
// If you have a goquery document just pass in doc.Selection.
func FromSelection(domain string, selec *goquery.Selection) string {
	opt := &Options{
		StrongDelimiter: "**",
	}
	markdown := SelecToMD(domain, selec, opt)

	return markdown
}

// FromString returns the content from a html string. If you
// already have a goquery selection use `FromSelection`.
func FromString(domain, html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", nil
	}
	return FromSelection(domain, doc.Selection), nil
}

// FromURL returns the content from the page with that url.
func FromURL(url string) (string, error) {
	// not using goquery.NewDocument directly because of the timeout
	resp, err := netClient.Get(url)
	if err != nil {
		return "", nil
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return "", nil
	}
	return FromSelection(DomainFromURL(url), doc.Selection), nil
}
