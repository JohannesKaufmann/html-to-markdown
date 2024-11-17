package cmd

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sebdah/goldie/v2"
)

func init() {
	OsExiter = func(code int) {
		// For the test cases we don't actually want to exit...
	}
}

type MockFileInfo struct {
	mode os.FileMode
}

func (info MockFileInfo) Name() string       { return "" }
func (info MockFileInfo) Size() int64        { return 1 }
func (info MockFileInfo) Mode() os.FileMode  { return info.mode }
func (info MockFileInfo) ModTime() time.Time { return time.Now() }
func (info MockFileInfo) IsDir() bool        { return false }
func (info MockFileInfo) Sys() interface{}   { return nil }

type FakeFile struct {
	bytes.Buffer
	mode os.FileMode
}

func (f FakeFile) Stat() (fs.FileInfo, error) {
	return &MockFileInfo{mode: f.mode}, nil
}

const (
	modePipe     = fs.FileMode(33554864) // "prw-rw----"
	modeTerminal = fs.FileMode(69206416) // "Dcrw--w----"
)

type CLIGoldenInput struct {
	modeStdin  os.FileMode
	modeStdout os.FileMode
	modeStderr os.FileMode

	inputStdin []byte
	inputArgs  []string
}

func cliGoldenTester(t *testing.T, input CLIGoldenInput) {
	if input.modeStdin == modeTerminal && input.inputStdin != nil {
		t.Fatal("invalid test: cannot provide stdin without pipe mode")
	}

	stdin := &FakeFile{mode: input.modeStdin}
	stdout := &FakeFile{mode: input.modeStdout}
	stderr := &FakeFile{mode: input.modeStderr}

	if input.inputStdin != nil {
		stdin.Write(input.inputStdin)
	}

	release := Release{
		Version: "2.3.4-test",
		Commit:  "ca82a6dff817ec66f44342007202690a93763949",
		Date:    "2024-08-18T13:03:43Z",
	}

	Run(stdin, stdout, stderr, input.inputArgs, release)

	if len(stdout.Bytes()) == 0 && len(stderr.Bytes()) == 0 {
		t.Fatal("neither stdout nor stderr have any content")
	}

	g := goldie.New(t)
	g.Assert(t, filepath.Join(t.Name(), "stdout"), stdout.Bytes())
	g.Assert(t, filepath.Join(t.Name(), "stderr"), stderr.Bytes())
}

type CLITestCase struct {
	desc string

	inputStdin []byte
	inputArgs  []string

	expectedStdout []byte
}

func cliSuccessTester(t *testing.T, tC CLITestCase) {
	stdin := &FakeFile{mode: modePipe}
	stdout := &FakeFile{mode: modePipe}
	stderr := &FakeFile{mode: modePipe}
	stdin.Write(tC.inputStdin)

	release := Release{
		Version: "2.3.4-test",
		Commit:  "ca82a6dff817ec66f44342007202690a93763949",
		Date:    "2024-08-18T13:03:43Z",
	}

	Run(stdin, stdout, stderr, tC.inputArgs, release)

	err := stderr.Bytes()
	if len(err) != 0 {
		t.Fatalf("got error: %q", string(err))
	}

	if !bytes.Equal(tC.expectedStdout, stdout.Bytes()) {
		t.Errorf("expected %q but got %q", string(tC.expectedStdout), stdout.String())
	}
}

func TestExecute(t *testing.T) {
	testCases := []struct {
		desc  string
		input CLIGoldenInput
	}{

		// - - - - - flag: version / help - - - - - //
		{
			desc: "[general] version terminal",

			input: CLIGoldenInput{
				modeStdin:  modeTerminal,
				modeStdout: modeTerminal,
				modeStderr: modeTerminal,

				inputArgs: []string{"html2markdown", "--version"},
			},
		},
		{
			desc: "[general] version pipe",

			input: CLIGoldenInput{
				modeStdin:  modePipe,
				modeStdout: modePipe,
				modeStderr: modePipe,

				inputArgs: []string{"html2markdown", "--version"},
			},
		},
		{
			desc: "[general] help terminal",

			input: CLIGoldenInput{
				modeStdin:  modeTerminal,
				modeStdout: modeTerminal,
				modeStderr: modeTerminal,

				inputArgs: []string{"html2markdown", "--help"},
			},
		},
		{
			desc: "[general] help pipe",

			input: CLIGoldenInput{
				modeStdin:  modePipe,
				modeStdout: modePipe,
				modeStderr: modePipe,

				inputArgs: []string{"html2markdown", "--help"},
			},
		},

		// - - - - - no content - - - - - //
		{
			desc: "[general] no content",

			input: CLIGoldenInput{
				modeStdin:  modeTerminal,
				modeStdout: modeTerminal,
				modeStderr: modeTerminal,

				inputStdin: nil,
				inputArgs:  []string{"html2markdown"},
			},
		},

		// - - - - - arguments - - - - - //
		{
			desc: "[argument unknown] version",

			input: CLIGoldenInput{
				modeStdin:  modeTerminal,
				modeStdout: modeTerminal,
				modeStderr: modeTerminal,

				inputArgs: []string{"html2markdown", `version`},
			},
		},
		{
			desc: "[argument unknown] html",

			input: CLIGoldenInput{
				modeStdin:  modeTerminal,
				modeStdout: modeTerminal,
				modeStderr: modeTerminal,

				inputArgs: []string{"html2markdown", `"<strong>text</strong>"`},
			},
		},
		{
			desc: "[argument unknown] long string",

			input: CLIGoldenInput{
				modeStdin:  modeTerminal,
				modeStdout: modeTerminal,
				modeStderr: modeTerminal,

				inputArgs: []string{"html2markdown", strings.Repeat("12456789", 40)},
			},
		},
		{
			desc: "[argument unknown] list of files",

			input: CLIGoldenInput{
				modeStdin:  modeTerminal,
				modeStdout: modeTerminal,
				modeStderr: modeTerminal,

				// The ** was treated as a file glob
				inputArgs: []string{"html2markdown", "--opt-strong-delimiter", "CONTRIBUTING.md", "README.md", "SECURITY.md", "a.html", "b.html", "c.html", "d.html", "e.html", "f.html"},
			},
		},

		// - - - - - flags - - - - - //
		{
			desc: "[flag unknown] with pipe",

			input: CLIGoldenInput{
				modeStdin:  modePipe,
				modeStdout: modePipe,
				modeStderr: modePipe,

				inputArgs: []string{"html2markdown", "--this-does-not-exist"},
			},
		},
		{
			desc: "[flag unknown] with terminal",

			input: CLIGoldenInput{
				modeStdin:  modeTerminal,
				modeStdout: modeTerminal,
				modeStderr: modeTerminal,

				inputArgs: []string{"html2markdown", "--this-does-not-exist"},
			},
		},

		{
			desc: "[flag misspelled] underscore",

			input: CLIGoldenInput{
				modeStdin:  modePipe,
				modeStdout: modePipe,
				modeStderr: modePipe,

				inputStdin: []byte("<strong>text</strong>"),
				// Someone accidentally used underscores instead of dashes
				inputArgs: []string{"html2markdown", "--opt_strong_delimiter="},
			},
		},

		// - - - - - converting - - - - - //
		{
			desc: "[convert] strong default",

			input: CLIGoldenInput{
				modeStdin:  modePipe,
				modeStdout: modePipe,
				modeStderr: modePipe,

				inputStdin: []byte("<strong>text</strong>"),
				inputArgs:  []string{"html2markdown"},
			},
		},
		{
			desc: "[convert] strong equal underscore",

			input: CLIGoldenInput{
				modeStdin:  modePipe,
				modeStdout: modePipe,
				modeStderr: modePipe,

				inputStdin: []byte("<strong>text</strong>"),
				// Note: We dont test the quoted version "__" since that is already unquoted by bash/go
				inputArgs: []string{"html2markdown", `--opt-strong-delimiter=__`},
			},
		},
		{
			desc: "[convert] strong space underscore",

			input: CLIGoldenInput{
				modeStdin:  modePipe,
				modeStdout: modePipe,
				modeStderr: modePipe,

				inputStdin: []byte("<strong>text</strong>"),
				inputArgs:  []string{"html2markdown", `--opt-strong-delimiter`, `__`},
			},
		},
		{
			desc: "[convert] collapse",

			input: CLIGoldenInput{
				modeStdin:  modePipe,
				modeStdout: modePipe,
				modeStderr: modePipe,

				inputStdin: []byte("<p>Some  <strong>   bold   </strong>  text</p>"),
				inputArgs:  []string{"html2markdown"},
			},
		},

		// - - - - - selectors - - - - - //
		{
			desc: "[include-selector] one match",

			input: CLIGoldenInput{
				modeStdin:  modePipe,
				modeStdout: modePipe,
				modeStderr: modePipe,

				inputStdin: []byte("<p>Some <strong><span>bold</span> text</strong> here</p>"),
				inputArgs:  []string{"html2markdown", "--include-selector", "strong"},
			},
		},
		{
			desc: "[include-selector] multiple matches",

			input: CLIGoldenInput{
				modeStdin:  modePipe,
				modeStdout: modePipe,
				modeStderr: modePipe,

				inputStdin: []byte("<p>Some <strong>a</strong> and <strong>b</strong> text</p>"),
				inputArgs:  []string{"html2markdown", "--include-selector", "strong"},
			},
		},
		{
			desc: "[include-selector] empty string",

			input: CLIGoldenInput{
				modeStdin:  modePipe,
				modeStdout: modePipe,
				modeStderr: modePipe,

				inputStdin: []byte("<p>Some <strong>a</strong> and <strong>b</strong> text</p>"),
				inputArgs:  []string{"html2markdown", "--include-selector", " "},
			},
		},
		{
			desc: "[include-selector] invalid",

			input: CLIGoldenInput{
				modeStdin:  modePipe,
				modeStdout: modePipe,
				modeStderr: modePipe,

				inputStdin: []byte("<p>Some <strong>a</strong> and <strong>b</strong> text</p>"),
				// This is not a valid selector, so cascadia is going to fail.
				inputArgs: []string{"html2markdown", "--include-selector", "?"},
			},
		},

		{
			desc: "[exclude-selector] exclude multiple",

			input: CLIGoldenInput{
				modeStdin:  modePipe,
				modeStdout: modePipe,
				modeStderr: modePipe,

				inputStdin: []byte(`<p>Some <strong>bold</strong> and <span class="italic">italic</span> text</p>`),
				inputArgs:  []string{"html2markdown", "--exclude-selector", "strong", "--exclude-selector", ".italic"},
			},
		},

		// - - - - - validation of options - - - - - //
		{
			desc: "[validation] no value",

			input: CLIGoldenInput{
				modeStdin:  modePipe,
				modeStdout: modePipe,
				modeStderr: modePipe,

				inputStdin: []byte("<strong>text</strong>"),
				inputArgs:  []string{"html2markdown", `--opt-strong-delimiter=`},
			},
		},
		{
			desc: "[validation] invalid value",

			input: CLIGoldenInput{
				modeStdin:  modePipe,
				modeStdout: modePipe,
				modeStderr: modePipe,

				inputStdin: []byte("<strong>text</strong>"),
				inputArgs:  []string{"html2markdown", `--opt-strong-delimiter=1234`},
			},
		},
		{
			desc: "[validation] discouraged value",

			input: CLIGoldenInput{
				modeStdin:  modePipe,
				modeStdout: modePipe,
				modeStderr: modePipe,

				inputStdin: []byte("<strong>text</strong>"),
				inputArgs:  []string{"html2markdown", `--opt-strong-delimiter=*`},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			cliGoldenTester(t, tC.input)
		})
	}
}

func TestExecute_General(t *testing.T) {
	testCases := []CLITestCase{
		{
			desc: "basic",

			inputStdin: []byte(`<p>Some <strong>a</strong> and <span class="bold">b</span> text</p>`),
			inputArgs:  []string{"html2markdown"},

			expectedStdout: []byte("Some **a** and b text\n"),
		},

		// - - - - - domain - - - - - //
		{
			desc: "[domain] without domain",

			inputStdin: []byte(`<img src="/image.png" />`),
			inputArgs:  []string{"html2markdown"},

			expectedStdout: []byte("![](/image.png)\n"),
		},
		{
			desc: "[domain] with domain",

			inputStdin: []byte(`<img src="/image.png" />`),
			inputArgs:  []string{"html2markdown", "--domain", "example.com"},

			expectedStdout: []byte("![](http://example.com/image.png)\n"),
		},
		{
			desc: "[domain] with https domain",

			inputStdin: []byte(`<img src="/image.png" />`),
			inputArgs:  []string{"html2markdown", "--domain", "https://example.com"},

			expectedStdout: []byte("![](https://example.com/image.png)\n"),
		},

		// - - - - - selectors - - - - - //
		{
			desc: "[include-selector] multiple matches",

			inputStdin: []byte(`<p>Some <strong>a</strong> and <span class="bold">b</span> text</p>`),
			inputArgs:  []string{"html2markdown", "--include-selector", "strong,.bold"},

			expectedStdout: []byte("**a**b\n"),
		},

		{
			desc: "[exclude-selector] exclude multiple with multiple flags",

			inputStdin: []byte(`<p>Some <strong>bold</strong> and <span class="italic">italic</span> text</p>`),
			inputArgs:  []string{"html2markdown", "--exclude-selector", "strong", "--exclude-selector", ".italic"},

			expectedStdout: []byte("Some and text\n"),
		},
		{
			desc: "[exclude-selector] exclude multiple with comma separator",

			inputStdin: []byte(`<p>Some <strong>bold</strong> and <span class="italic">italic</span> text</p>`),
			inputArgs:  []string{"html2markdown", "--exclude-selector", "strong,.italic"},

			expectedStdout: []byte("Some and text\n"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			cliSuccessTester(t, tC)
		})
	}
}

func TestExecute_Plugins(t *testing.T) {
	testCases := []CLITestCase{

		{
			desc: "[plugin-strikethrough] disabled by default",

			inputStdin: []byte(`<p>Some <s>outdated</s> text</p>`),
			inputArgs:  []string{"html2markdown"},

			expectedStdout: []byte("Some outdated text\n"),
		},
		{
			desc: "[plugin-strikethrough] enabled",

			inputStdin: []byte(`<p>Some <s>outdated</s> text</p>`),
			inputArgs:  []string{"html2markdown", "--plugin-strikethrough"},

			expectedStdout: []byte("Some ~~outdated~~ text\n"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			cliSuccessTester(t, tC)
		})
	}
}
