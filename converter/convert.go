package converter

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"

	"golang.org/x/net/html"
)

type convertOption struct {
	domain  string
	context context.Context
}
type convertOptionFunc func(o *convertOption)

func WithContext(ctx context.Context) convertOptionFunc {
	return func(o *convertOption) {
		o.context = ctx
	}
}
func WithDomain(domain string) convertOptionFunc {
	return func(o *convertOption) {
		o.domain = domain
	}
}

func (conv *Converter) setError(err error) {
	conv.m.Lock()
	defer conv.m.Unlock()

	conv.err = err
}
func (conv *Converter) getError() error {
	conv.m.RLock()
	defer conv.m.RUnlock()

	return conv.err
}

var errNoRenderHandlers = errors.New("no render handlers are registered. did you forget to register the commonmark plugin?")

// ConvertNode converts a `*html.Node` to a markdown byte slice.
//
// If you have already parsed an HTML page using the `html.Parse()` function
// from the "golang.org/x/net/html" package then you can pass this node
// directly to the converter.
func (conv *Converter) ConvertNode(doc *html.Node, opts ...convertOptionFunc) ([]byte, error) {

	if err := conv.getError(); err != nil {
		// There can be errors while calling `Init` on the plugins (e.g. validation errors).
		// Now is the first opportunity where we can return that error.
		return nil, err
	}

	conv.m.Lock()
	option := &convertOption{}
	for _, fn := range opts {
		fn(option)
	}
	conv.m.Unlock()

	// If there are no render handlers registered this is
	// usually a user error - since people want the Commonmark Plugin in 99% of cases.
	if len(conv.getRenderHandlers()) == 0 {
		// TODO: Add Name() to the interface & check for the presence of *both* the Base & Commonmark Plugin
		// TODO: What if just the base plugin is registered?
		return nil, errNoRenderHandlers
	}

	// - - - - - - - - - - - - - - - - - - - //

	state := newGlobalState()

	if option.context == nil {
		option.context = context.Background()
	}
	ctx := option.context
	ctx = provideDomain(ctx, option.domain)
	ctx = provideAssembleAbsoluteURL(ctx, defaultAssembleAbsoluteURL)
	ctx = state.provideGlobalState(ctx)

	customCtx := newConverterContext(ctx, conv)

	// - - - - - - - - - - - - - - - - - - - //

	// Pre-Render
	for _, handler := range conv.getPreRenderHandlers() {
		handler.Value(customCtx, doc)
	}

	// Render
	var buf bytes.Buffer
	conv.handleRenderNode(customCtx, &buf, doc)

	// Post-Render
	result := buf.Bytes()
	for _, handler := range conv.getPostRenderHandlers() {
		result = handler.Value(customCtx, result)
	}

	return result, nil
}

// ConvertReader converts the html from the reader to markdown.
//
// Under the hood `html.Parse()` is used to parse the HTML.
func (conv *Converter) ConvertReader(r io.Reader, opts ...convertOptionFunc) ([]byte, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	return conv.ConvertNode(doc, opts...)
}

// ConvertString converts a html-string to a markdown-string.
//
// Under the hood `html.Parse()` is used to parse the HTML.
func (conv *Converter) ConvertString(htmlInput string, opts ...convertOptionFunc) (string, error) {
	r := strings.NewReader(htmlInput)
	output, err := conv.ConvertReader(r, opts...)
	if err != nil {
		return "", err
	}

	return string(output), nil
}
