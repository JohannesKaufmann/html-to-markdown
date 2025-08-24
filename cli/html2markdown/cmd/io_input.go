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
	inputFullFilepath  string
	outputFullFilepath string
	data               []byte
}

// E.g. "website.html" -> "website"
func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

// If the output is a file, it would be "output.md"
var defaultBasename = "output"

func (cli *CLI) listInputs() ([]*input, error) {

	// NOTE: When both stdin and --input are specified,
	// the explicit --file argument takes precedence.
	// This improves interoperability with other tools like `xargs`.
	// https://github.com/JohannesKaufmann/html-to-markdown/issues/170
	if cli.config.inputFilepath != "" {
		matches, err := doublestar.FilepathGlob(cli.config.inputFilepath, doublestar.WithFilesOnly(), doublestar.WithNoFollow())
		if err != nil {
			return nil, err
		}
		if len(matches) == 0 {
			// The inputFilepath wasn't actually a glob but was pointing to an existing folder.
			// The user probably wanted to convert all files in that folder — so we recommend the glob.
			if outInfo, err := os.Stat(cli.config.inputFilepath); err == nil && outInfo.IsDir() {
				return nil, NewCLIError(
					fmt.Errorf("input path %q is a directory, not a file", cli.config.inputFilepath),
					Paragraph("Here is how you can use a glob to match multiple files:"),
					CodeBlock(`html2markdown --input "src/*.html" --output "dist/"`),
				)
			}

			return nil, NewCLIError(
				fmt.Errorf("no files found matching pattern %q", cli.config.inputFilepath),
				Paragraph("Here is how you can use a glob to match multiple files:"),
				CodeBlock(`html2markdown --input "src/*.html" --output "dist/"`),
			)
		}

		var inputs []*input
		for _, match := range matches {
			inputs = append(inputs, &input{
				inputFullFilepath: match,
				data:              nil,
			})
		}

		return inputs, nil
	}

	if cli.isStdinPipe {
		data, err := io.ReadAll(cli.Stdin)
		if err != nil {
			return nil, err
		}
		return []*input{
			{
				inputFullFilepath: defaultBasename,
				data:              data,
			},
		}, nil
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

	data, err := os.ReadFile(in.inputFullFilepath)
	if err != nil {
		return nil, err
	}

	return data, nil
}
