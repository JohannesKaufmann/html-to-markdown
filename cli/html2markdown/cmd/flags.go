package cmd

import (
	"errors"
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

func (cli *CLI) singleStringFlag(target *string, name string, usage string) {
	cli.flags.Func(name, usage, func(flagValue string) error {
		if strings.TrimSpace(flagValue) == "" {
			return errors.New("empty string")
		}
		if *target != "" {
			return fmt.Errorf("another value has already been set")
		}

		*target = strings.TrimSpace(flagValue)
		return nil
	})
}

func (cli *CLI) initFlags(progname string) {
	cli.flags = flag.NewFlagSet(progname, flag.ContinueOnError)
	cli.flags.SetOutput(io.Discard)

	// - - - - - General - - - - - //
	cli.flags.BoolVar(&cli.config.version, "version", false, "display the version")
	cli.flags.BoolVar(&cli.config.version, "v", false, "display the version")

	cli.singleStringFlag(
		&cli.config.inputFilepath,
		"input",
		"Read input from FILE instead of stdin",
	)
	cli.singleStringFlag(
		&cli.config.outputFilepath,
		"output",
		"Write output to FILE instead of stdout",
	)
	cli.flags.BoolVar(&cli.config.outputOverwrite, "output-overwrite", false, "replace existing files")

	// TODO: --tag-type-block=script,style (and check that it is not a selector)
	// TODO: --tag-type-inline=script,style (and check that it is not a selector)

	cli.flags.StringVar(
		&cli.config.domain,
		"domain",
		"",
		"The url of the web page, used to convert relative links to absolute links.",
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

	cli.flags.BoolVar(&cli.config.enablePluginTable, "plugin-table", false, "enable the plugin table")
	cli.flags.BoolVar(&cli.config.tableSkipEmptyRows, "opt-table-skip-empty-rows", false, "[for --plugin-table] omit empty rows from the output")
	cli.flags.BoolVar(&cli.config.tableHeaderPromotion, "opt-table-header-promotion", false, "[for --plugin-table] first row should be treated as a header")
	cli.flags.StringVar(&cli.config.tableSpanCellBehavior, "opt-table-span-cell-behavior", "", `[for --plugin-table] how colspan/rowspan should be rendered: "empty" or "mirror"`)
	cli.flags.BoolVar(&cli.config.tablePresentationTables, "opt-table-presentation-tables", false, `[for --plugin-table] whether tables with role="presentation" should be converted`)
	cli.flags.StringVar(&cli.config.tableNewlineBehavior, "opt-table-newline-behavior", "", `[for --plugin-table] how tables containing newlines should be handled: "skip" or "preserve"`)
	cli.flags.BoolVar(&cli.config.tablePadColumns, "opt-table-pad-columns", true, `[for --plugin-table] whether columns in the tables should include extra padding for visual continuity`)
}

func (cli *CLI) parseFlags(args []string) error {
	err := cli.flags.Parse(args)
	if err != nil {
		return cli.categorizeFlagError(err)
	}

	cli.config.args = cli.flags.Args()

	// Validate flag dependencies
	if cli.config.tableSkipEmptyRows && !cli.config.enablePluginTable {
		return fmt.Errorf("--opt-table-skip-empty-rows requires --plugin-table to be enabled")
	}
	if cli.config.tableHeaderPromotion && !cli.config.enablePluginTable {
		return fmt.Errorf("--opt-table-header-promotion requires --plugin-table to be enabled")
	}
	if cli.config.tableSpanCellBehavior != "" && !cli.config.enablePluginTable {
		return fmt.Errorf("--opt-table-span-cell-behavior requires --plugin-table to be enabled")
	}
	if cli.config.tablePresentationTables && !cli.config.enablePluginTable {
		return fmt.Errorf("--opt-table-presentation-tables requires --plugin-table to be enabled")
	}
	if cli.config.tableNewlineBehavior != "" && !cli.config.enablePluginTable {
		return fmt.Errorf("--opt-table-newline-behavior requires --plugin-table to be enabled")
	}
	if cli.config.tablePadColumns && !cli.config.enablePluginTable {
		return fmt.Errorf("--opt-table-pad-columns requires --plugin-table to be enabled")
	}

	// TODO: use constant for flag name & use formatFlag
	//       var keyStrongDelimiter = "opt-strong-delimiter"

	return nil
}
