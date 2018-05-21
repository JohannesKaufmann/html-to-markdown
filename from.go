package md

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type RuleFunc func(content string, selec *goquery.Selection, options *Options) *string

type Converter struct {
	sync.RWMutex
	rules map[string][]RuleFunc

	domain  string
	options Options
}

func NewConverter(domain string, enableCommonmark bool, options *Options) *Converter {
	c := &Converter{
		domain: domain,
		rules:  make(map[string][]RuleFunc),
	}

	if enableCommonmark {
		c.AddRules(commonmark...)
	}

	// TODO: put domain in options?
	if options == nil {
		options = &Options{}
	}
	if options.StrongDelimiter == "" {
		options.StrongDelimiter = "**"
	}
	if options.Fence == "" {
		options.Fence = "```"
	}
	if options.HR == "" {
		options.HR = "* * *"
	} else {
		fmt.Println("count *", strings.Count(options.HR, "*"))
		// validateOptions()
	}

	c.options = *options
	// ...

	return c
}
func (c *Converter) getRuleFuncs(tag string) []RuleFunc {
	c.RLock()
	defer c.RUnlock()

	r, ok := c.rules[tag]
	if !ok || len(r) == 0 {
		return []RuleFunc{DefaultRule}
	}

	return r
}
func (c *Converter) AddRules(rules ...Rule) *Converter {
	c.Lock()
	defer c.Unlock()

	for _, rule := range rules {
		for _, filter := range rule.Filter {
			r, _ := c.rules[filter]
			r = append(r, rule.Replacement)
			c.rules[filter] = r
		}
	}

	return c
}

// TODO: naming -> Run, Proccess, ...

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

var multipleNewLinesRegex = regexp.MustCompile(`[\n]{2,}`)

// Convert returns the content from a goquery selection.
// If you have a goquery document just pass in doc.Selection.
func (c *Converter) Convert(selec *goquery.Selection) string {
	c.RLock() // DONT NEED THIS?
	domain := c.domain
	options := c.options
	l := len(c.rules)
	if l == 0 {
		panic("you have added no rules. either enable commonmark or add you own.")
	}
	c.RUnlock()

	markdown := c.selecToMD(domain, selec, &options)

	markdown = strings.TrimSpace(markdown)
	markdown = multipleNewLinesRegex.ReplaceAllString(markdown, "\n\n")

	return markdown
}

func (c *Converter) ConvertReader(reader io.Reader) (bytes.Buffer, error) {
	var buffer bytes.Buffer
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return buffer, err
	}
	buffer.WriteString(
		c.Convert(doc.Selection),
	)

	return buffer, nil
}
func (c *Converter) ConvertResponse(res *http.Response) (string, error) {
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return "", err
	}
	return c.Convert(doc.Selection), nil
}

// ConvertString returns the content from a html string. If you
// already have a goquery selection use `Convert`.
func (c *Converter) ConvertString(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}
	return c.Convert(doc.Selection), nil
}

func (c *Converter) ConvertBytes(bytes []byte) ([]byte, error) {
	res, err := c.ConvertString(string(bytes))
	if err != nil {
		return nil, err
	}
	return []byte(res), nil
}

// ConvertURL returns the content from the page with that url.
func (c *Converter) ConvertURL(url string) (string, error) {
	// not using goquery.NewDocument directly because of the timeout
	resp, err := netClient.Get(url)
	if err != nil {
		return "", nil
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return "", nil
	}
	domain := DomainFromURL(url)
	if c.domain != domain {
		return "", errors.New("expected " + c.domain + " as the domain but got " + domain)
	}
	return c.Convert(doc.Selection), nil
}
