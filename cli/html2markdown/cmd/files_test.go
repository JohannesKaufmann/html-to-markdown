package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newTestDir(t *testing.T) string {
	t.Helper()
	var tempFolderPattern = "html2markdown_*_testdata"

	directoryPath, err := os.MkdirTemp("", tempFolderPattern)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("creating a temp directory: %q\n", directoryPath)

	return directoryPath
}
func newTestDirWithFiles(t *testing.T) string {
	t.Helper()
	directoryPath := newTestDir(t)

	inputFolder := filepath.Join(directoryPath, "input")
	nestedFolder := filepath.Join(directoryPath, "input", "nested")
	outputFolder := filepath.Join(directoryPath, "output")

	// - - -  /input/  - - - //
	err := os.MkdirAll(inputFolder, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(filepath.Join(inputFolder, "random.txt"), []byte("other random file"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(filepath.Join(inputFolder, "website_a.html"), []byte("<strong>file content A</strong>"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(filepath.Join(inputFolder, "website_b.html"), []byte("<strong>file content B</strong>"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// - - -  /input/nested/  - - - //
	err = os.MkdirAll(nestedFolder, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(filepath.Join(nestedFolder, "website_c.html"), []byte("<i>file content C</i>"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// - - -  /output/  - - - //
	err = os.MkdirAll(outputFolder, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	return directoryPath
}

func writePipeChar(buf *bytes.Buffer, index int) {
	if index == 0 {
		return
	}
	buf.WriteString(strings.Repeat("│ ", index-1))
	buf.WriteString("├─")
}
func writeFileRepresentation(buf *bytes.Buffer, rel string, info os.FileInfo, data []byte) {
	parts := strings.Split(rel, string(os.PathSeparator))

	if rel == "." {
		writePipeChar(buf, 0)
	} else {
		writePipeChar(buf, len(parts))
	}

	buf.WriteString(filepath.Base(rel))

	if info.IsDir() {
		// no extra info
	} else {
		buf.WriteString(fmt.Sprintf(" %q", string(data)))
	}
}

// Similar to render representation from the "dom" package.
func renderRepresentation(dirPath string) (string, error) {
	var buf bytes.Buffer

	err := filepath.Walk(dirPath, func(name string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		var rel string
		rel, err = filepath.Rel(dirPath, name)
		if err != nil {
			return err
		}

		if info.IsDir() {
			writeFileRepresentation(&buf, rel, info, nil)
		} else {
			data, err := os.ReadFile(name)
			if err != nil {
				return err
			}
			writeFileRepresentation(&buf, rel, info, data)
		}
		buf.WriteRune('\n')

		return nil
	})
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func expectRepresentation(t *testing.T, directoryPath string, expectedFS string) {
	actualFS, err := renderRepresentation(directoryPath)
	if err != nil {
		t.Fatal(err)
	}

	actualFS = strings.TrimSpace(actualFS)
	expectedFS = strings.TrimSpace(expectedFS)

	if actualFS != expectedFS {
		t.Errorf("expected \n%s\nbut got\n%s", expectedFS, actualFS)
	}
}

var testRelease = Release{
	Version: "2.3.4-test",
	Commit:  "ca82a6dff817ec66f44342007202690a93763949",
	Date:    "2024-08-18T13:03:43Z",
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - //

func TestExecute_SingleFileInput(t *testing.T) {
	directoryPath := newTestDir(t)
	defer os.RemoveAll(directoryPath)

	expectedStdout := []byte("**the file content**\n")
	inputPath := filepath.Join(directoryPath, "input.html")

	err := os.WriteFile(inputPath, []byte("<strong>the file content</strong>"), 0644)
	if err != nil {
		t.Error(err)
	}

	// - - - - - - - - - //
	args := []string{"html2markdown", "--input", inputPath}

	stdin := &FakeFile{mode: modeTerminal}
	stdout := &FakeFile{mode: modePipe}
	stderr := &FakeFile{mode: modePipe}

	Run(stdin, stdout, stderr, args, testRelease)

	stderrBytes := stderr.Bytes()
	if len(stderrBytes) != 0 {
		t.Fatalf("got error: %q", string(stderrBytes))
	}
	if !bytes.Equal(expectedStdout, stdout.Bytes()) {
		t.Errorf("expected %q but got %q", string(expectedStdout), stdout.String())
	}
	// - - - - - - - - - //

	expectRepresentation(t, directoryPath, `
.
├─input.html "<strong>the file content</strong>"
	`)
}

func TestExecute_SingleFileOutput(t *testing.T) {
	directoryPath := newTestDir(t)
	defer os.RemoveAll(directoryPath)

	outputPath := directoryPath + string(filepath.Separator)

	// - - - - - - - - - //
	args := []string{"html2markdown", "--output", outputPath}

	stdin := &FakeFile{mode: modePipe}
	stdout := &FakeFile{mode: modePipe}
	stderr := &FakeFile{mode: modePipe}
	stdin.WriteString("<strong>bold text</strong>")

	Run(stdin, stdout, stderr, args, testRelease)

	stderrBytes := stderr.Bytes()
	if len(stderrBytes) != 0 {
		t.Fatalf("got error: %q", string(stderrBytes))
	}
	stdoutBytes := stdout.Bytes()
	if len(stdoutBytes) != 0 {
		t.Fatalf("got content: %q", string(stdoutBytes))
	}
	// - - - - - - - - - //

	expectRepresentation(t, directoryPath, `
.
├─output.md "**bold text**"
	`)
}

func TestExecute_DirectoryOutput(t *testing.T) {
	directoryPath := newTestDir(t)
	defer os.RemoveAll(directoryPath)

	inputPath := filepath.Join(directoryPath, "my_website.html")

	err := os.WriteFile(inputPath, []byte("<strong>the file content</strong>"), 0644)
	if err != nil {
		t.Error(err)
	}

	// - - - - - - - - - //
	args := []string{"html2markdown", "--input", inputPath, "--output", directoryPath + string(os.PathSeparator)}

	stdin := &FakeFile{mode: modeTerminal}
	stdout := &FakeFile{mode: modePipe}
	stderr := &FakeFile{mode: modePipe}

	Run(stdin, stdout, stderr, args, testRelease)

	if len(stderr.Bytes()) != 0 {
		t.Fatalf("expected no stderr content but got %q", stderr.String())
	}
	if len(stdout.Bytes()) != 0 {
		t.Fatalf("expected no stdout content")
	}
	// - - - - - - - - - //

	expectRepresentation(t, directoryPath, `
.
├─my_website.html "<strong>the file content</strong>"
├─my_website.md "**the file content**"
	`)
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - //

func TestExecute_FilePattern(t *testing.T) {
	t.Run("the default test dir with files", func(t *testing.T) {
		directoryPath := newTestDirWithFiles(t)
		defer os.RemoveAll(directoryPath)

		// This is the default file structure *before*
		// the CLI does any changes:
		expectRepresentation(t, directoryPath, `
.
├─input
│ ├─nested
│ │ ├─website_c.html "<i>file content C</i>"
│ ├─random.txt "other random file"
│ ├─website_a.html "<strong>file content A</strong>"
│ ├─website_b.html "<strong>file content B</strong>"
├─output
		`)
	})

	testCases := []struct {
		desc         string
		assembleArgs func(dir string) []string
		expected     string
	}{
		{
			desc: "output to a specific file",
			assembleArgs: func(dir string) []string {
				input := filepath.Join(dir, "input", "website_a.html")
				output := filepath.Join(dir, "output", "websites", "the_cool_website.md")

				return []string{"html2markdown", "--input", input, "--output", output}
			},
			expected: `
.
├─input
│ ├─nested
│ │ ├─website_c.html "<i>file content C</i>"
│ ├─random.txt "other random file"
│ ├─website_a.html "<strong>file content A</strong>"
│ ├─website_b.html "<strong>file content B</strong>"
├─output
│ ├─websites
│ │ ├─the_cool_website.md "**file content A**"

			`,
		},
		{
			desc: "override a specific file",
			assembleArgs: func(dir string) []string {
				input := filepath.Join(dir, "input", "website_a.html")
				output := filepath.Join(dir, "input", "website_a.html")

				return []string{"html2markdown", "--input", input, "--output", output}
			},
			// Note: The file content of "website_a.html" changed
			expected: `
.
├─input
│ ├─nested
│ │ ├─website_c.html "<i>file content C</i>"
│ ├─random.txt "other random file"
│ ├─website_a.html "**file content A**"
│ ├─website_b.html "<strong>file content B</strong>"
├─output

			`,
		},
		// - - - - - - - - - - - - pattern - - - - - - - - - - - - //
		{
			desc: "pattern matches single file",
			assembleArgs: func(dir string) []string {
				input := filepath.Join(dir, "input", "website_a.*")
				output := filepath.Join(dir, "output")

				return []string{"html2markdown", "--input", input, "--output", output + string(os.PathSeparator)}
			},
			expected: `
.
├─input
│ ├─nested
│ │ ├─website_c.html "<i>file content C</i>"
│ ├─random.txt "other random file"
│ ├─website_a.html "<strong>file content A</strong>"
│ ├─website_b.html "<strong>file content B</strong>"
├─output
│ ├─website_a.md "**file content A**"
			`,
		},
		{
			desc: "direct website html files",
			assembleArgs: func(dir string) []string {
				input := filepath.Join(dir, "input", "website*.html")
				output := filepath.Join(dir, "output")

				return []string{"html2markdown", "--input", input, "--output", output + string(os.PathSeparator)}
			},
			expected: `
.
├─input
│ ├─nested
│ │ ├─website_c.html "<i>file content C</i>"
│ ├─random.txt "other random file"
│ ├─website_a.html "<strong>file content A</strong>"
│ ├─website_b.html "<strong>file content B</strong>"
├─output
│ ├─website_a.md "**file content A**"
│ ├─website_b.md "**file content B**"
			`,
		},
		{
			desc: "match everything",
			assembleArgs: func(dir string) []string {
				input := filepath.Join(dir, "**", "*")
				output := filepath.Join(dir, "output")

				return []string{"html2markdown", "--input", input, "--output", output + string(os.PathSeparator)}
			},

			// Note: The "random.md" was also placed in the output folder.
			expected: `
.
├─input
│ ├─nested
│ │ ├─website_c.html "<i>file content C</i>"
│ ├─random.txt "other random file"
│ ├─website_a.html "<strong>file content A</strong>"
│ ├─website_b.html "<strong>file content B</strong>"
├─output
│ ├─random.md "other random file"
│ ├─website_a.md "**file content A**"
│ ├─website_b.md "**file content B**"
│ ├─website_c.md "*file content C*"
			`,
		},
		{
			desc: "nested website html files",
			assembleArgs: func(dir string) []string {
				input := filepath.Join(dir, "**", "website*.html")
				output := filepath.Join(dir, "output", "in", "nested", "folder")

				return []string{"html2markdown", "--input", input, "--output", output + string(os.PathSeparator)}
			},
			expected: `
.
├─input
│ ├─nested
│ │ ├─website_c.html "<i>file content C</i>"
│ ├─random.txt "other random file"
│ ├─website_a.html "<strong>file content A</strong>"
│ ├─website_b.html "<strong>file content B</strong>"
├─output
│ ├─in
│ │ ├─nested
│ │ │ ├─folder
│ │ │ │ ├─website_a.md "**file content A**"
│ │ │ │ ├─website_b.md "**file content B**"
│ │ │ │ ├─website_c.md "*file content C*"
			`,
		},
		{
			desc: "input and output is same folder",
			assembleArgs: func(dir string) []string {
				input := filepath.Join(dir, "input", "website*.html")
				output := filepath.Join(dir, "input")

				return []string{"html2markdown", "--input", input, "--output", output + string(os.PathSeparator)}
			},
			expected: `
.
├─input
│ ├─nested
│ │ ├─website_c.html "<i>file content C</i>"
│ ├─random.txt "other random file"
│ ├─website_a.html "<strong>file content A</strong>"
│ ├─website_a.md "**file content A**"
│ ├─website_b.html "<strong>file content B</strong>"
│ ├─website_b.md "**file content B**"
├─output
			`,
		},
		{
			desc: "input and output is same folder and nested pattern",
			assembleArgs: func(dir string) []string {
				input := filepath.Join(dir, "**", "website*.html")
				output := filepath.Join(dir, "input")

				return []string{"html2markdown", "--input", input, "--output", output + string(os.PathSeparator)}
			},

			// Note: "website_c.md" was placed in "input" because by default the flat output structure is used.
			expected: `
.
├─input
│ ├─nested
│ │ ├─website_c.html "<i>file content C</i>"
│ ├─random.txt "other random file"
│ ├─website_a.html "<strong>file content A</strong>"
│ ├─website_a.md "**file content A**"
│ ├─website_b.html "<strong>file content B</strong>"
│ ├─website_b.md "**file content B**"
│ ├─website_c.md "*file content C*"
├─output
			`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			directoryPath := newTestDirWithFiles(t)
			defer os.RemoveAll(directoryPath)

			args := tC.assembleArgs(directoryPath)
			t.Logf("args: %+v\n", args)

			// - - - - - - - - - //
			stdin := &FakeFile{mode: modeTerminal}
			stdout := &FakeFile{mode: modePipe}
			stderr := &FakeFile{mode: modePipe}

			Run(stdin, stdout, stderr, args, testRelease)

			if len(stderr.Bytes()) != 0 {
				t.Fatalf("expected no stderr content but got %q", stderr.String())
			}
			if len(stdout.Bytes()) != 0 {
				t.Fatalf("expected no stdout content")
			}
			// - - - - - - - - - //

			expectRepresentation(t, directoryPath, tC.expected)
		})
	}
}
