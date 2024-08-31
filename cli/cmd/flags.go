package cmd

import (
	"flag"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type FlagString string

func (a *FlagString) Scan(state fmt.ScanState, verb rune) error {
	token, err := state.Token(true, func(r rune) bool {
		return unicode.IsLetter(r) || r == '-'
	})
	if err != nil {
		return err
	}
	*a = FlagString(token)
	return nil
}

func flagStringSlice(elems *[]string) func(string) error {
	return func(raw string) error {
		values := strings.Split(raw, ",")

		for _, val := range values {
			val = strings.TrimSpace(val)
			if val == "" {
				continue
			}

			*elems = append(*elems, val)
		}
		return nil
	}
}

func (cli *CLI) initFlags(progname string) {
	cli.flags = flag.NewFlagSet(progname, flag.ContinueOnError)
	cli.flags.SetOutput(io.Discard)

	// - - - //

	cli.flags.BoolVar(&cli.config.version, "version", false, "display the version")
	cli.flags.BoolVar(&cli.config.version, "v", false, "display the version")

	// cli.flags.BoolVar(&cli.config.help, "help", false, "display help")

	cli.flags.StringVar(
		&cli.config.strongDelimiter,
		"opt-strong-delimiter",
		"**",
		`Make bold text. Should <strong> be indicated by two asterisks or two underscores?
"**" or "__" (default: "**")`,
	)

	// cli.flags.StringVar(&cli.config.strongDelimiter, "opt-heading-style", "", "")
	// cli.flags.StringVar(&cli.config.strongDelimiter, "opt-horizontal-rule", "", "")
	// cli.flags.StringVar(&cli.config.strongDelimiter, "opt-bullet-list-marker", "", "")

	// TODO: how to disable commonmark plugin?
	// --plugin_commonmark=false
	// --plugin.commonmark=false
	// --no-plugin="cm" / --disable-plugin="cm"
	// But what if we have conflicting flags???
	cli.flags.Func("plugins", "which plugins should be enabled?", flagStringSlice(&cli.config.plugins))
}

func (cli *CLI) parseFlags(args []string) error {
	err := cli.flags.Parse(args)
	if err != nil {
		return cli.categorizeFlagError(err)
	}

	cli.config.args = cli.flags.Args()

	return nil
}
