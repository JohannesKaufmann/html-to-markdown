package cmd

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
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
			fmt.Errorf(`when processing multiple input files, --output "%s" must end with "%s" to indicate a directory`, dir, dir+"/"),
		)
	}

	// TODO: The glob can also be a folder with just one file...
	//       So we should check if the path contains any glob characters.

	// Check if output path exists
	if outInfo, err := os.Stat(outputPath); err == nil {
		if outInfo.IsDir() {
			dir := filepath.Base(outputPath)
			return "", NewCLIError(
				fmt.Errorf(`path "%s" exists and is a directory, did you mean "%s"?`, dir, dir+"/"),
				Paragraph(`The --output must end with "/" to indicate a directory`),
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

func calculateOutputPaths(inputFilepath string, inputs []*input) error {
	globBase, _ := doublestar.SplitPattern(inputFilepath)

	allBasenames := make(map[string]int)
	for _, input := range inputs {
		basenameWithExt := filepath.Base(input.inputFullFilepath)
		basename := fileNameWithoutExtension(basenameWithExt)

		val := allBasenames[basename]
		if val == 0 {
			// -> The standard filename
			input.outputFullFilepath = basename + ".md"
		} else {
			relativePath, err := filepath.Rel(globBase, input.inputFullFilepath)
			if err != nil {
				return err
			}

			fmt.Printf("input:%q globBase:%q relativePath:%q\n", input.inputFullFilepath, globBase, relativePath)

			relativePath2, err := filepath.Rel(filepath.FromSlash(globBase), input.inputFullFilepath)
			if err != nil {
				return err
			}
			fmt.Printf("relativePath2:%q\n", relativePath2)

			// We hash the relative path (based from the globBase)
			// since the globBase is *the same* for all files.
			// Bonus: It makes testing easier as the temporary folder does not matter.
			hash := hashFilepath(relativePath)

			// -> The filename for duplicates
			input.outputFullFilepath = basename + "." + hash[:10] + ".md"
		}

		allBasenames[basename]++
	}

	return nil
}

func hashFilepath(path string) string {
	h := sha256.New()
	h.Write([]byte(
		filepath.ToSlash(path),
	))

	bs := h.Sum(nil)

	fmt.Printf("hashFilepath %q -> %q -> %q\n", path, filepath.ToSlash(path), fmt.Sprintf("%x", bs))

	return fmt.Sprintf("%x", bs)
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
			err := WriteFile(filepath.Join(cli.config.outputFilepath, filename), markdown, cli.config.outputOverwrite)
			if err != nil {
				if errors.Is(err, os.ErrExist) {
					return fmt.Errorf("output path %q already exists. Use --output-overwrite to replace existing files", cli.config.outputFilepath)
				}

				return fmt.Errorf("error while writing the file into the directory: %w", err)
			}

			return nil
		}
	case outputTypeFile:
		{
			err := WriteFile(cli.config.outputFilepath, markdown, cli.config.outputOverwrite)
			if err != nil {
				if errors.Is(err, os.ErrExist) {
					return fmt.Errorf("output path %q already exists. Use --output-overwrite to replace existing files", cli.config.outputFilepath)
				}

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

// WriteFile writes data to a file with override control
// If override is false and file exists, returns an error
// If override is true, truncates existing file or creates new one
func WriteFile(filename string, data []byte, override bool) error {
	// As the base flags we have:
	//   O_WRONLY = write to the file, not read
	//   O_CREATE = create the file if it doesn't exist
	flag := os.O_WRONLY | os.O_CREATE

	if override {
		// We add this flag:
		//   O_TRUNC = the existing contents are truncated to zero length
		flag |= os.O_TRUNC
	} else {
		// We add this flag:
		//   O_EXCL = if used with O_CREATE, causes error if file already exists
		flag |= os.O_EXCL
	}

	f, err := os.OpenFile(filename, flag, 0644)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}
	return err
}
