package converter

import (
	"github.com/JohannesKaufmann/dom"
	"golang.org/x/net/html"
)

type register struct {
	conv *Converter
}

// - - - - - - - - - - - - - Pre-Render - - - - - - - - - - - - - //

type HandlePreRenderFunc func(ctx Context, doc *html.Node)

func (r *register) PreRenderer(fn HandlePreRenderFunc, priority int) {
	r.conv.m.Lock()
	defer r.conv.m.Unlock()

	handler := prioritized(fn, priority)
	r.conv.preRenderHandlers = append(r.conv.preRenderHandlers, handler)
}
func (conv *Converter) getPreRenderHandlers() prioritizedSlice[HandlePreRenderFunc] {
	conv.m.RLock()
	defer conv.m.RUnlock()

	handlers := make(prioritizedSlice[HandlePreRenderFunc], len(conv.preRenderHandlers))
	copy(handlers, conv.preRenderHandlers)
	handlers.Sort()

	return handlers
}

// - - - - - - - - - - - - - Render - - - - - - - - - - - - - //

// Writer is an interface that only conforms to the Write* methods of bytes.Buffer
type Writer interface {
	Write(p []byte) (n int, err error)
	WriteByte(c byte) error
	WriteRune(r rune) (n int, err error)
	WriteString(s string) (n int, err error)
}

type HandleRenderFunc func(ctx Context, w Writer, n *html.Node) RenderStatus

func (r *register) Renderer(fn HandleRenderFunc, priority int) {
	r.conv.m.Lock()
	defer r.conv.m.Unlock()

	handler := prioritized(fn, priority)
	r.conv.renderHandlers = append(r.conv.renderHandlers, handler)
}
func (conv *Converter) getRenderHandlers() prioritizedSlice[HandleRenderFunc] {
	conv.m.RLock()
	defer conv.m.RUnlock()

	handlers := make(prioritizedSlice[HandleRenderFunc], len(conv.renderHandlers))
	copy(handlers, conv.renderHandlers)
	handlers.Sort()

	return handlers
}

// - - - - - - - - - - - - - Post Render - - - - - - - - - - - - - //

type HandlePostRenderFunc func(ctx Context, content []byte) []byte

func (r *register) PostRenderer(fn HandlePostRenderFunc, priority int) {
	r.conv.m.Lock()
	defer r.conv.m.Unlock()

	handler := prioritized(fn, priority)
	r.conv.postRenderHandlers = append(r.conv.postRenderHandlers, handler)
}
func (conv *Converter) getPostRenderHandlers() prioritizedSlice[HandlePostRenderFunc] {
	conv.m.RLock()
	defer conv.m.RUnlock()

	handlers := make(prioritizedSlice[HandlePostRenderFunc], len(conv.postRenderHandlers))
	copy(handlers, conv.postRenderHandlers)
	handlers.Sort()

	return handlers
}

// - - - - - - - - - - - - - Text - - - - - - - - - - - - - //

type HandleTextTransformFunc func(ctx Context, content string) string

func (r *register) TextTransformer(fn HandleTextTransformFunc, priority int) {
	r.conv.m.Lock()
	defer r.conv.m.Unlock()

	handler := prioritized(fn, priority)
	r.conv.textTransformHandlers = append(r.conv.textTransformHandlers, handler)
}
func (conv *Converter) getTextTransformHandlers() prioritizedSlice[HandleTextTransformFunc] {
	conv.m.RLock()
	defer conv.m.RUnlock()

	handlers := make(prioritizedSlice[HandleTextTransformFunc], len(conv.textTransformHandlers))
	copy(handlers, conv.textTransformHandlers)
	handlers.Sort()

	return handlers
}

// - - - - - - - - - - - - - Escaping - - - - - - - - - - - - - //

func (r *register) EscapedChar(chars ...rune) {
	r.conv.m.Lock()
	defer r.conv.m.Unlock()

	for _, char := range chars {
		r.conv.markdownChars[char] = struct{}{}
	}
}
func (conv *Converter) checkIsEscapedChar(r rune) bool {
	conv.m.RLock()
	defer conv.m.RUnlock()

	_, ok := conv.markdownChars[r]
	return ok
}

type HandleUnEscapeFunc func(chars []byte, index int) int

func (r *register) UnEscaper(fn HandleUnEscapeFunc, priority int) {
	r.conv.m.Lock()
	defer r.conv.m.Unlock()

	handler := prioritized(fn, priority)
	r.conv.unEscapeHandlers = append(r.conv.unEscapeHandlers, handler)
}
func (conv *Converter) getUnEscapeHandlers() prioritizedSlice[HandleUnEscapeFunc] {
	conv.m.RLock()
	defer conv.m.RUnlock()

	handlers := make(prioritizedSlice[HandleUnEscapeFunc], len(conv.unEscapeHandlers))
	copy(handlers, conv.unEscapeHandlers)
	handlers.Sort()

	return handlers
}

// - - - - - - - - - - - - - Tag Strategy - - - - - - - - - - - - - //

func (r *register) TagStrategy(tagName string, strategy tagStrategy) {
	r.conv.m.Lock()
	defer r.conv.m.Unlock()

	r.conv.tagStrategies[tagName] = strategy
}
func (conv *Converter) getTagStrategy(tagName string) (tagStrategy, bool) {
	conv.m.RLock()
	defer conv.m.RUnlock()

	strategy, ok := conv.tagStrategies[tagName]
	return strategy, ok
}
func (conv *Converter) getTagStrategyWithFallback(tagName string) tagStrategy {
	decision, ok := conv.getTagStrategy(tagName)
	if ok {
		return decision
	}

	if dom.NameIsBlockNode(tagName) {
		return StrategyMarkdownBlock
	}
	return StrategyMarkdownLeaf
}
