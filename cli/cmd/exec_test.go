package cmd

import (
	"bytes"
	"fmt"
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
		fmt.Println("OS_EXITER_CALLED", code)
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

type CLIInput struct {
	modeStdin  os.FileMode
	modeStdout os.FileMode
	modeStderr os.FileMode

	inputStdin []byte
	inputArgs  []string
}

func cliTester(t *testing.T, input CLIInput) {
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

func TestExecute(t *testing.T) {
	testCases := []struct {
		desc  string
		input CLIInput
	}{

		// - - - - - flag: version / help - - - - - //
		{
			desc: "[general] version terminal",

			input: CLIInput{
				modeStdin:  modeTerminal,
				modeStdout: modeTerminal,
				modeStderr: modeTerminal,

				inputArgs: []string{"html2markdown", "--version"},
			},
		},
		{
			desc: "[general] version pipe",

			input: CLIInput{
				modeStdin:  modePipe,
				modeStdout: modePipe,
				modeStderr: modePipe,

				inputArgs: []string{"html2markdown", "--version"},
			},
		},
		{
			desc: "[general] help terminal",

			input: CLIInput{
				modeStdin:  modeTerminal,
				modeStdout: modeTerminal,
				modeStderr: modeTerminal,

				inputArgs: []string{"html2markdown", "--help"},
			},
		},
		{
			desc: "[general] help pipe",

			input: CLIInput{
				modeStdin:  modePipe,
				modeStdout: modePipe,
				modeStderr: modePipe,

				inputArgs: []string{"html2markdown", "--help"},
			},
		},

		// - - - - - no content - - - - - //
		{
			desc: "[general] no content",

			input: CLIInput{
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

			input: CLIInput{
				modeStdin:  modeTerminal,
				modeStdout: modeTerminal,
				modeStderr: modeTerminal,

				inputArgs: []string{"html2markdown", `version`},
			},
		},
		{
			desc: "[argument unknown] html",

			input: CLIInput{
				modeStdin:  modeTerminal,
				modeStdout: modeTerminal,
				modeStderr: modeTerminal,

				inputArgs: []string{"html2markdown", `"<strong>text</strong>"`},
			},
		},
		{
			desc: "[argument unknown] long string",

			input: CLIInput{
				modeStdin:  modeTerminal,
				modeStdout: modeTerminal,
				modeStderr: modeTerminal,

				inputArgs: []string{"html2markdown", strings.Repeat("12456789", 40)},
			},
		},
		{
			desc: "[argument unknown] list of files",

			input: CLIInput{
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

			input: CLIInput{
				modeStdin:  modePipe,
				modeStdout: modePipe,
				modeStderr: modePipe,

				inputArgs: []string{"html2markdown", "--this-does-not-exist"},
			},
		},
		{
			desc: "[flag unknown] with terminal",

			input: CLIInput{
				modeStdin:  modeTerminal,
				modeStdout: modeTerminal,
				modeStderr: modeTerminal,

				inputArgs: []string{"html2markdown", "--this-does-not-exist"},
			},
		},

		{
			desc: "[flag misspelled] underscore",

			input: CLIInput{
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

			input: CLIInput{
				modeStdin:  modePipe,
				modeStdout: modePipe,
				modeStderr: modePipe,

				inputStdin: []byte("<strong>text</strong>"),
				inputArgs:  []string{"html2markdown"},
			},
		},
		{
			desc: "[convert] strong equal underscore",

			input: CLIInput{
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

			input: CLIInput{
				modeStdin:  modePipe,
				modeStdout: modePipe,
				modeStderr: modePipe,

				inputStdin: []byte("<strong>text</strong>"),
				inputArgs:  []string{"html2markdown", `--opt-strong-delimiter`, `__`},
			},
		},
		{
			desc: "[convert] collapse",

			input: CLIInput{
				modeStdin:  modePipe,
				modeStdout: modePipe,
				modeStderr: modePipe,

				inputStdin: []byte("<p>Some  <strong>   bold   </strong>  text</p>"),
				inputArgs:  []string{"html2markdown"},
			},
		},

		// - - - - - validation of options - - - - - //
		{
			desc: "[validation] no value",

			input: CLIInput{
				modeStdin:  modePipe,
				modeStdout: modePipe,
				modeStderr: modePipe,

				inputStdin: []byte("<strong>text</strong>"),
				inputArgs:  []string{"html2markdown", `--opt-strong-delimiter=`},
			},
		},
		{
			desc: "[validation] invalid value",

			input: CLIInput{
				modeStdin:  modePipe,
				modeStdout: modePipe,
				modeStderr: modePipe,

				inputStdin: []byte("<strong>text</strong>"),
				inputArgs:  []string{"html2markdown", `--opt-strong-delimiter=1234`},
			},
		},
		{
			desc: "[validation] discouraged value",

			input: CLIInput{
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
			cliTester(t, tC.input)
		})
	}
}
