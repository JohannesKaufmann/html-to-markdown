package converter

import "sync"

type Converter struct {
	m sync.RWMutex

	err error

	preRenderHandlers  prioritizedSlice[HandlePreRenderFunc]
	renderHandlers     prioritizedSlice[HandleRenderFunc]
	postRenderHandlers prioritizedSlice[HandlePostRenderFunc]

	textTransformHandlers prioritizedSlice[HandleTextTransformFunc]

	markdownChars    map[rune]interface{}
	unEscapeHandlers prioritizedSlice[HandleUnEscapeFunc]

	tagStrategies map[string]tagStrategy

	Register register
}

type converterOption = func(c *Converter) error

func NewConverter(opts ...converterOption) *Converter {
	conv := &Converter{
		markdownChars: make(map[rune]interface{}),
		tagStrategies: make(map[string]tagStrategy),
	}
	conv.Register = register{conv}

	conv.registerBase()

	for _, opt := range opts {
		err := opt(conv)
		if err != nil {
			conv.setError(err)
			break
		}
	}

	return conv
}
