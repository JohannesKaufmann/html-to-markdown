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

	// If there are no render handlers registered (apart from the base) this is
	// usually a user error - since people want the Commonmark Plugin in 99% of cases.
	countBaseRenderHandlers := 1
	if len(conv.getRenderHandlers()) == countBaseRenderHandlers {
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

func (conv *Converter) ConvertReader(r io.Reader, opts ...convertOptionFunc) ([]byte, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	return conv.ConvertNode(doc, opts...)
}

func (conv *Converter) ConvertString(htmlInput string, opts ...convertOptionFunc) (string, error) {
	r := strings.NewReader(htmlInput)
	output, err := conv.ConvertReader(r, opts...)
	if err != nil {
		return "", err
	}

	return string(output), nil
}
