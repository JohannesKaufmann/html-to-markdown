package cmd

import (
	"fmt"
	"io"

	"github.com/muesli/termenv"
)

type CLIError struct {
	cause    error
	printers []Printer
}

func extractCLIError(err error) (CLIError, bool) {
	if cliErr, ok := err.(*CLIError); ok {
		return *cliErr, true
	}

	return CLIError{
		cause: err,
	}, false
}

func NewCLIError(cause error, printers ...Printer) error {
	return &CLIError{
		cause:    cause,
		printers: printers,
	}
}
func (e CLIError) Error() string {
	return e.cause.Error()
}
func (e CLIError) PrintDetails(w io.Writer) {
	errPrinter := ColoredBox("error", e.cause.Error())

	// Prepend the error printer
	e.printers = append([]Printer{errPrinter}, e.printers...)

	for _, printer := range e.printers {
		w.Write([]byte("\n"))
		printer.Print(w)
	}
	w.Write([]byte("\n"))
}

func (cli CLI) PrintErr(err error) {
	if err == nil {
		return
	}

	e, _ := extractCLIError(err)
	e.PrintDetails(cli.Stderr)
}
func (cli CLI) PrintWarn(err error) {
	if err == nil {
		return
	}

	output := termenv.NewOutput(cli.Stderr)

	prefix := output.String("warning:").Background(termenv.ANSIYellow).Foreground(termenv.ANSIBrightWhite).String()
	message := output.String(err.Error()).Foreground(termenv.ANSIYellow).String()

	fmt.Fprintf(cli.Stderr, "\n%s %s\n\n", prefix, message)
}
