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

// chdirWithCleanup changes the current working directory to the named directory,
// and then restore the original working directory at the end of the test.
//
// TODO: Once we are on 1.24 we can replace this with t.Chdir()
func chdirWithCleanup(t *testing.T, dir string) {
	olddir, err := os.Getwd()
	if err != nil {
		t.Fatalf("chdir: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir %s: %v", dir, err)
	}

	t.Cleanup(func() {
		if err := os.Chdir(olddir); err != nil {
			t.Errorf("chdir to original working directory %s: %v", olddir, err)
			os.Exit(1)
		}
	})
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
	pathSeparator := "/" // <-- we don't use os.PathSeparator here just to test that windows also supports slash
	args := []string{"html2markdown", "--input", inputPath, "--output", directoryPath + pathSeparator}

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

func TestExecute_NotOverwrite(t *testing.T) {
	directoryPath := newTestDir(t)
	defer os.RemoveAll(directoryPath)

	outputPath := filepath.Join(directoryPath, "output.md")

	t.Run("the first run", func(t *testing.T) {
		args := []string{"html2markdown", "--output", outputPath}

		stdin := &FakeFile{mode: modePipe}
		stdout := &FakeFile{mode: modePipe}
		stderr := &FakeFile{mode: modePipe}
		stdin.WriteString("<strong>file content A</strong>")

		Run(stdin, stdout, stderr, args, testRelease)

		stderrBytes := stderr.Bytes()
		if len(stderrBytes) != 0 {
			t.Fatalf("got error: %q", string(stderrBytes))
		}
		if len(stdout.Bytes()) != 0 {
			t.Fatalf("expected no stdout content")
		}

		expectRepresentation(t, directoryPath, `
.
├─output.md "**file content A**"
		`)
	})

	t.Run("the second run", func(t *testing.T) {
		args := []string{"html2markdown", "--output", outputPath}

		stdin := &FakeFile{mode: modePipe}
		stdout := &FakeFile{mode: modePipe}
		stderr := &FakeFile{mode: modePipe}
		stdin.WriteString("<strong>file content B</strong>")

		Run(stdin, stdout, stderr, args, testRelease)

		actualStderr := stderr.String()
		expectedStderr := fmt.Sprintf("\nerror: output path %q already exists. Use --output-overwrite to replace existing files\n\n", outputPath)
		if actualStderr != expectedStderr {
			t.Errorf("expected stderr %q but got %q", expectedStderr, actualStderr)
		}
		if len(stdout.Bytes()) != 0 {
			t.Fatalf("expected no stdout content")
		}

		// The content should still be the same:
		expectRepresentation(t, directoryPath, `
.
├─output.md "**file content A**"
		`)
	})
}

func TestExecute_DuplicateFiles(t *testing.T) {
	directoryPath := newTestDir(t)
	defer os.RemoveAll(directoryPath)

	chdirWithCleanup(t, directoryPath)

	inputFolderA := filepath.Join(directoryPath, "input", "a")
	inputFolderB := filepath.Join(directoryPath, "input", "b")
	inputFolderC := filepath.Join(directoryPath, "input", "nested", "c")

	err := os.MkdirAll(inputFolderA, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(inputFolderB, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(inputFolderC, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(filepath.Join(inputFolderA, "random.html"), []byte("file a"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(filepath.Join(inputFolderB, "random.html"), []byte("file b"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(filepath.Join(inputFolderC, "random.html"), []byte("file c"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// - - - - - - - - - //
	args := []string{"html2markdown", "--input", filepath.Join(".", "input", "**", "*"), "--output", filepath.Join(".", "output") + "/"}

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
├─input
│ ├─a
│ │ ├─random.html "file a"
│ ├─b
│ │ ├─random.html "file b"
│ ├─nested
│ │ ├─c
│ │ │ ├─random.html "file c"
├─output
│ ├─random.689330a60f.md "file b"
│ ├─random.f679b6e0c2.md "file c"
│ ├─random.md "file a"
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
			desc: "output to a specific extension",
			assembleArgs: func(dir string) []string {
				input := filepath.Join(dir, "input", "website_a.html")
				output := filepath.Join(dir, "output", "websites", "the_cool_website.txt")

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
│ │ ├─the_cool_website.txt "**file content A**"
			`,
		},
		{
			desc: "override a specific file",
			assembleArgs: func(dir string) []string {
				input := filepath.Join(dir, "input", "website_a.html")
				output := filepath.Join(dir, "input", "website_a.html")

				return []string{"html2markdown", "--input", input, "--output", output, "--output-overwrite"}
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

		// - - - - - - - - - - - - relative path - - - - - - - - - - - - //
		{
			desc: "relative path: output to a specific file",
			assembleArgs: func(_ string) []string {
				input := filepath.Join(".", "input", "website_a.html")
				output := filepath.Join(".", "output", "websites", "the_cool_website.md")

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
			desc: "relative path: direct website html files",
			assembleArgs: func(_ string) []string {
				input := filepath.Join(".", "input", "website*.html")
				output := filepath.Join(".", "output") + string(os.PathSeparator)

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
│ ├─website_a.md "**file content A**"
│ ├─website_b.md "**file content B**"
			`,
		},
		{
			desc: "relative path: output to current directory",
			assembleArgs: func(_ string) []string {
				input := filepath.Join(".", "input", "website*.html")
				output := "." + string(os.PathSeparator)

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
├─website_a.md "**file content A**"
├─website_b.md "**file content B**"
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

			chdirWithCleanup(t, directoryPath)

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

func TestWriteFile_OverrideFalse(t *testing.T) {
	directoryPath := newTestDir(t)
	defer os.RemoveAll(directoryPath)

	filePath := filepath.Join(directoryPath, "test.txt")
	override := false

	expectRepresentation(t, directoryPath, ".")

	// - - - file does not exist yet - - - //
	err := WriteFile(filePath, []byte("A"), override)
	if err != nil {
		t.Error(err)
	}
	expectRepresentation(t, directoryPath, ".\n"+`├─test.txt "A"`)

	// - - - file exists already - - - //
	err = WriteFile(filePath, []byte("B"), override)
	if err == nil {
		t.Error("expected there to be an error but got nil")
	}
	expectRepresentation(t, directoryPath, ".\n"+`├─test.txt "A"`) // <-- still the old content
}
func TestWriteFile_OverrideTrue(t *testing.T) {
	directoryPath := newTestDir(t)
	defer os.RemoveAll(directoryPath)

	filePath := filepath.Join(directoryPath, "test.txt")
	override := true

	expectRepresentation(t, directoryPath, ".")

	// - - - file does not exist yet - - - //
	err := WriteFile(filePath, []byte("A"), override)
	if err != nil {
		t.Error(err)
	}
	expectRepresentation(t, directoryPath, ".\n"+`├─test.txt "A"`)

	// - - - file exists already - - - //
	err = WriteFile(filePath, []byte("B"), override)
	if err != nil {
		t.Error(err)
	}
	expectRepresentation(t, directoryPath, ".\n"+`├─test.txt "B"`) // <-- the new content
}
