package cmd

import (
	"bytes"
	"fmt"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
)

func overrideValidationError(e *commonmark.ValidateConfigError) error {

	// TODO: Maybe OptionFunc should already validate and return an error?
	//       Then it would be easier to override the Key since we have once
	//       place to assemble the []OptionFunc and directly treat the errors...
	//
	// We would basically invoke it ourselves:
	//    err := commonmark.WithStrongDelimiter(cli.config.strongDelimiter)(conv)

	switch e.Key {
	case "StrongDelimiter":
		e.Key = "opt-strong-delimiter"
	}

	e.KeyWithValue = fmt.Sprintf("--%s=%q", e.Key, e.Value)
	return e
}
func (cli *CLI) convert(input []byte) ([]error, error) {

	conv := converter.NewConverter(
		converter.WithPlugins(
			// TODO: there should be a golden file test for functionality that needs the base plugin
			base.NewBasePlugin(),
			commonmark.NewCommonmarkPlugin(
				commonmark.WithStrongDelimiter(cli.config.strongDelimiter),
			),
		),
	)

	r := bytes.NewReader(input)
	markdown, err := conv.ConvertReader(r)
	if err != nil {
		e, ok := err.(*commonmark.ValidateConfigError)
		if ok {
			return nil, overrideValidationError(e)
		}

		return nil, err
	}

	fmt.Fprintln(cli.Stdout, string(markdown))
	return nil, nil
}
