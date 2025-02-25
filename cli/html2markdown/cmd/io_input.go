package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

type input struct {
	fullFilepath string
	data         []byte
}

// E.g. "website.html" -> "website"
func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func getOutputFileName(fullFilepath string) string {
	basenameWithExt := filepath.Base(fullFilepath)
	basename := fileNameWithoutExtension(basenameWithExt)

	return basename + ".md"
}

func (cli *CLI) listInputs() ([]*input, error) {
	if cli.isStdinPipe && cli.config.inputFilepath != "" {
		return nil, NewCLIError(
			fmt.Errorf("cannot use both stdin and --input at the same time. Use either stdin or specify an input file, but not both"),
		)
	}

	if cli.isStdinPipe {
		data, err := io.ReadAll(cli.Stdin)
		if err != nil {
			return nil, err
		}
		return []*input{
			{
				// If the output is a file, it would be "output.md"
				fullFilepath: "output",
				data:         data,
			},
		}, nil
	}

	if cli.config.inputFilepath != "" {
		matches, err := doublestar.FilepathGlob(cli.config.inputFilepath, doublestar.WithFilesOnly(), doublestar.WithNoFollow())
		if err != nil {
			return nil, err
		}
		if len(matches) == 0 {
			return nil, NewCLIError(
				fmt.Errorf("no files found matching pattern %q", cli.config.inputFilepath),
				Paragraph("Here is how you can use a glob to match multiple files:"),
				CodeBlock(`html2markdown --input "src/*.html" --output "dist/"`),
			)
		}

		var inputs []*input
		for _, match := range matches {
			inputs = append(inputs, &input{
				fullFilepath: match,
				data:         nil,
			})
		}

		return inputs, nil
	}

	return nil, NewCLIError(
		fmt.Errorf("the html input should be piped into the cli"),
		Paragraph("Here is how you can use the CLI:"),
		CodeBlock(`echo "<strong>important</strong>" | html2markdown`),
	)
}

func (cli *CLI) readInput(in *input) ([]byte, error) {
	if in.data != nil {
		return in.data, nil
	}

	data, err := os.ReadFile(in.fullFilepath)
	if err != nil {
		return nil, err
	}

	return data, nil
}
