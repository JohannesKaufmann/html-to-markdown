package cmd

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/andybalholm/cascadia"
)

var (
	projectBinary = "html2markdown"
)

// OsExiter is the function used when the app exits. If not set defaults to os.Exit.
var OsExiter = os.Exit

// - - - - - - - - - - - - - //

type Config struct {
	// args are the positional (non-flag) command-line arguments.
	args []string

	inputFilepath   string
	outputFilepath  string
	outputOverwrite bool

	// - - - - - General - - - - - //
	version bool
	domain  string

	includeSelector cascadia.SelectorGroup
	excludeSelector cascadia.SelectorGroup

	// - - - - - Options - - - - - //
	strongDelimiter string

	// - - - - - Plugins - - - - - //
	enablePluginStrikethrough bool

	enablePluginTable       bool
	tableSkipEmptyRows      bool
	tableHeaderPromotion    bool
	tableSpanCellBehavior   string
	tablePresentationTables bool
	tableNewlineBehavior    string
}

// Release holds the information (from the 3 ldflags) that goreleaser sets.
type Release struct {
	// Current Git tag (the v prefix is stripped)
	Version string

	// Current git commit SHA
	Commit string

	// Date in the RFC3339 format
	Date string
}
type CLI struct {
	Stdin  ReadWriterWithStat
	Stdout ReadWriterWithStat
	Stderr ReadWriterWithStat

	OsArgs []string

	Release Release

	isStdinPipe  bool
	isStdoutPipe bool
	isStderrPipe bool

	flags  *flag.FlagSet
	config Config

	usageText bytes.Buffer
}

func (cli *CLI) Init() error {
	var err error
	cli.isStdinPipe, err = isPipe(cli.Stdin)
	if err != nil {
		return fmt.Errorf("error while checking stdin for is pipe: %w", err)
	}
	cli.isStdoutPipe, err = isPipe(cli.Stdout)
	if err != nil {
		return fmt.Errorf("error while checking stdout for is pipe: %w", err)
	}
	cli.isStderrPipe, err = isPipe(cli.Stderr)
	if err != nil {
		return fmt.Errorf("error while checking stderr for is pipe: %w", err)
	}

	cli.initFlags(cli.OsArgs[0])

	err = cli.initUsageText()
	if err != nil {
		return fmt.Errorf("error while initializing the usage text: %w", err)
	}

	return nil
}
func (cli *CLI) Execute() {

	warnings, err := cli.run()

	for _, warning := range warnings {
		cli.PrintWarn(warning)
	}

	if err == flag.ErrHelp {
		cli.printUsage()

		OsExiter(0)
		return
	} else if err != nil {
		cli.PrintErr(err)

		OsExiter(1) // General Error
		return
	} else {
		OsExiter(0)
		return
	}
}

func (cli *CLI) run() ([]error, error) {

	err := cli.parseFlags(cli.OsArgs[1:])
	if err != nil {
		return nil, err
	}

	if len(cli.config.args) != 0 {

		return nil, NewCLIError(
			fmt.Errorf("unknown arguments: %s", strings.Join(cli.config.args, " ")),
			Paragraph("Here is how you can use the CLI:"),
			CodeBlock(`echo "<strong>important</strong>" | html2markdown`),
		)
	}

	if cli.config.version {
		cli.printVersion()
		return nil, nil
	}

	// - - - - - - - - - - - - - - - //

	inputs, err := cli.listInputs()
	if err != nil {
		return nil, err
	}

	outputType, err := determineOutputType(cli.config.inputFilepath, len(inputs), cli.config.outputFilepath)
	if err != nil {
		return nil, err
	}

	err = ensureOutputDirectories(outputType, cli.config.outputFilepath)
	if err != nil {
		return nil, err
	}

	err = calculateOutputPaths(cli.config.inputFilepath, inputs)
	if err != nil {
		return nil, err
	}

	for _, input := range inputs {
		data, err := cli.readInput(input)
		if err != nil {
			return nil, err
		}

		markdown, err := cli.convert(data)
		if err != nil {
			return nil, err
		}

		err = cli.writeOutput(outputType, input.outputFullFilepath, markdown)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}
