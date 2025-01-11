package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

type file struct {
	name      string
	extension string

	input []byte
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func (cli *CLI) getInputs() ([]*file, error) {
	if cli.isStdinPipe {
		data, err := io.ReadAll(cli.Stdin)
		if err != nil {
			return nil, err
		}
		return []*file{
			{
				name:  "output",
				input: data,
			},
		}, nil
	}

	if cli.config.inputFilepath != "" {
		matches, err := doublestar.FilepathGlob(cli.config.inputFilepath)
		if err != nil {
			return nil, err
		}
		if len(matches) == 0 {
			return nil, NewCLIError(
				fmt.Errorf("no files found matching pattern %q", cli.config.inputFilepath),
				Paragraph("Here is how you can use a glob to match multiple files:"),
				CodeBlock(`html2markdown --input "src/*.html" --output "dist"`),
			)
		}
		// if len(matches) != 1 {
		// 	return nil, errors.New("converting multiple files at once is not (yet) supported")
		// }

		var files []*file
		for _, match := range matches {
			data, err := os.ReadFile(match)
			if err != nil {
				return nil, err
			}

			filename := filepath.Base(match)

			fmt.Printf("BASE:%q EXT:%q NAME:%q \n", filepath.Base(match), filepath.Ext(match), fileNameWithoutExtension(filename))

			files = append(files, &file{
				name:      fileNameWithoutExtension(filename),
				extension: filepath.Ext(filename),

				input: data,
			})
		}

		return files, nil
	}

	return nil, NewCLIError(
		fmt.Errorf("the html input should be piped into the cli"),
		Paragraph("Here is how you can use the CLI:"),
		CodeBlock(`echo "<strong>important</strong>" | html2markdown`),
	)
}
