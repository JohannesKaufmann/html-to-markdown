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

func determineOutputType(countInputs int, outputPath string) outputType {
	if outputPath == "" {
		return outputTypeStdout
	}

	// TODO: the glob can also be a folder with just one file...
	if countInputs > 1 {
		// There are multiple inputs, so the input MUST have been a glob or directory.
		// It also means that the output MUST be a directory.
		return outputTypeDirectory
	}

	if strings.HasSuffix(outputPath, "/") || strings.HasSuffix(outputPath, "\\") {
		// With the trailing slash a directory can be indicated.
		return outputTypeDirectory
	}

	if filepath.Ext(filepath.Base(outputPath)) != "" {
		// With a file extension it is LIKELY to be a file.
		return outputTypeFile
	}

	// Check if output path exists
	if outInfo, err := os.Stat(outputPath); err == nil {
		if outInfo.IsDir() {
			return outputTypeDirectory
		}
		return outputTypeFile
	}

	// Default to file for single input
	return outputTypeFile
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
