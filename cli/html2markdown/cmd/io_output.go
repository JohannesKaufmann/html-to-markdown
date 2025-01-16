package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type outputType string

const (
	outputTypeStdout    outputType = "stdout"
	outputTypeDirectory outputType = "directory"
	outputTypeFile      outputType = "file"
)

// The user can indicate that they mean a directory by having a slash as the suffix.
func hasFolderSuffix(outputPath string) bool {
	// Note: We generally support the os.PathSeparator (e.g. "\" on windows).
	//       But also "/" is always supported.
	return strings.HasSuffix(outputPath, string(os.PathSeparator)) || strings.HasSuffix(outputPath, "/")
}

func determineOutputType(_inputPath string, countInputs int, outputPath string) (outputType, error) {
	if outputPath == "" {
		if countInputs > 1 {
			return "", NewCLIError(
				fmt.Errorf("when processing multiple input files --output needs to be a directory"),
				Paragraph("Here is how you can use a glob to match multiple files:"),
				CodeBlock(`html2markdown --input "src/*.html" --output "dist/"`),
			)
		}

		return outputTypeStdout, nil
	}

	if hasFolderSuffix(outputPath) {
		return outputTypeDirectory, nil
	}

	// - - - - - - - - - //
	// We can now assume that the output path specifies a file.
	// But let's make sure...

	if countInputs > 1 {
		// There are multiple inputs, so the input MUST have been a glob or directory.
		// It also means that the output MUST be a directory.
		dir := filepath.Base(outputPath)
		return "", NewCLIError(
			fmt.Errorf(`when processing multiple input files, --output "%s" must end with "%s" to indicate a directory`, dir, dir+string(os.PathSeparator)),
		)
	}

	// TODO: The glob can also be a folder with just one file...
	//       So we should check if the path contains any glob characters.

	// Check if output path exists
	if outInfo, err := os.Stat(outputPath); err == nil {
		if outInfo.IsDir() {
			dir := filepath.Base(outputPath)
			return "", NewCLIError(
				fmt.Errorf(`path "%s" exists and is a directory - did you mean "%s"?`, dir, dir+string(os.PathSeparator)),
				Paragraph(fmt.Sprintf(`The --output must end with "%s" to indicate a directory`, string(os.PathSeparator))),
			)
		}
		return outputTypeFile, nil
	}

	if filepath.Ext(filepath.Base(outputPath)) != "" {
		// With a file extension it is LIKELY to be a file.
		return outputTypeFile, nil
	}

	// Default to file for single input
	return outputTypeFile, nil
}

func ensureOutputDirectories(outputType outputType, outputFilepath string) error {
	if outputType == outputTypeDirectory {

		return os.MkdirAll(outputFilepath, os.ModePerm)
	} else if outputType == outputTypeFile {
		path := filepath.Dir(outputFilepath)
		return os.MkdirAll(path, os.ModePerm)
	} else {
		return nil
	}
}

func (cli *CLI) writeOutput(outputType outputType, filename string, markdown []byte) error {
	switch outputType {
	case outputTypeDirectory:
		{
			err := os.WriteFile(filepath.Join(cli.config.outputFilepath, filename), markdown, 0644)
			if err != nil {
				return fmt.Errorf("error while writing the file into the directory: %w", err)
			}

			return nil
		}
	case outputTypeFile:
		{
			err := os.WriteFile(cli.config.outputFilepath, markdown, 0644)
			if err != nil {
				return fmt.Errorf("error while writing the file: %w", err)
			}

			return nil
		}
	default:
		{
			fmt.Fprintln(cli.Stdout, string(markdown))

			return nil
		}
	}
}
