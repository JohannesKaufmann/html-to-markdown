package cmd

import (
	"fmt"
	"io"

	"github.com/muesli/termenv"
)

type Printer interface {
	Print(w io.Writer)
}

// - - - - - - - //

type coloredBox struct {
	prefix string
	text   string
}

func ColoredBox(prefix string, text string) Printer {
	return &coloredBox{prefix, text}
}

func (p coloredBox) Print(w io.Writer) {
	output := termenv.NewOutput(w)

	prefix := output.String(p.prefix + ":").Background(termenv.ANSIRed).Foreground(termenv.ANSIBrightWhite).String()
	message := output.String(p.text).Foreground(termenv.ANSIRed).String()

	fmt.Fprintf(w, "%s %s\n", prefix, message)
}

// - - - - - - - //

type paragraph struct {
	text string
}

func Paragraph(text string) Printer {
	return &paragraph{text}
}
func (p paragraph) Print(w io.Writer) {
	fmt.Fprintln(w, p.text)
}

// - - - - - - - //

type codeBlock struct {
	code string
}

func CodeBlock(code string) Printer {
	return &codeBlock{code}
}
func (cb codeBlock) Print(w io.Writer) {
	// TODO: what about indenting multiline?
	fmt.Fprintf(w, "    %s\n", cb.code)
}
