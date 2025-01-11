package cmd

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/andybalholm/cascadia"
)

// selectorFlag sets up a flag that parses a CSS selector string into a cascadia.Selector.
func (cli *CLI) selectorFlag(target *cascadia.SelectorGroup, name string, usage string) {
	cli.flags.Func(name, usage, func(flagValue string) error {
		if strings.TrimSpace(flagValue) == "" {
			return fmt.Errorf("invalid css selector: empty string")
		}

		// Compile the provided CSS selector string
		sel, err := cascadia.ParseGroup(flagValue)
		if err != nil {
			return fmt.Errorf("invalid css selector: %w", err)
		}

		*target = append(*target, sel...)
		return nil
	})
}

func (cli *CLI) initFlags(progname string) {
	cli.flags = flag.NewFlagSet(progname, flag.ContinueOnError)
	cli.flags.SetOutput(io.Discard)

	// - - - - - General - - - - - //
	cli.flags.BoolVar(&cli.config.version, "version", false, "display the version")
	cli.flags.BoolVar(&cli.config.version, "v", false, "display the version")

	cli.flags.StringVar(
		&cli.config.inputFilepath,
		"input",
		"",
		"Read input from FILE instead of stdin",
	)
	cli.flags.StringVar(
		&cli.config.outputFilepath,
		"output",
		"",
		"Write output to FILE instead of stdout",
	)

	// TODO: --tag-type-block=script,style (and check that it is not a selector)
	// TODO: --tag-type-inline=script,style (and check that it is not a selector)

	cli.flags.StringVar(
		&cli.config.domain,
		"domain",
		"",
		"domain of the web page, needed for links",
	)

	cli.selectorFlag(&cli.config.includeSelector, "include-selector", "css query selector to only include parts of the input")
	cli.selectorFlag(&cli.config.excludeSelector, "exclude-selector", "css query selector to exclude parts of the input")

	// - - - - - Options - - - - - //
	cli.flags.StringVar(
		&cli.config.strongDelimiter,
		"opt-strong-delimiter",
		"**",
		`Make bold text. Should <strong> be indicated by two asterisks or two underscores?
"**" or "__" (default: "**")`,
	)

	// - - - - - Plugins - - - - - //
	// TODO: --opt-strikethrough-delimiter for the strikethrough plugin
	cli.flags.BoolVar(&cli.config.enablePluginStrikethrough, "plugin-strikethrough", false, "enable the plugin ~~strikethrough~~")

}

func (cli *CLI) parseFlags(args []string) error {
	err := cli.flags.Parse(args)
	if err != nil {
		return cli.categorizeFlagError(err)
	}

	cli.config.args = cli.flags.Args()

	return nil
}
