package tester

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sebdah/goldie/v2"
)

var enableRoundTrip = flag.Bool("round", false, "enable the round trip testing")

const suffixInputFile = ".in.html"
const suffixOutputFile = ".out.md"

func getInputFiles(pathOfFolder string) ([]string, error) {
	files, err := os.ReadDir(pathOfFolder)
	if err != nil {
		return nil, fmt.Errorf("error while reading %q folder: %w", pathOfFolder, err)
	}

	var names []string
	for _, file := range files {
		if file.IsDir() {
			return nil, fmt.Errorf("did not expected a folder %q", file.Name())
		}
		if strings.HasSuffix(file.Name(), suffixOutputFile) {
			continue
		}
		if !strings.HasSuffix(file.Name(), suffixInputFile) {
			return nil, fmt.Errorf("only expect in or out files but got %q", file.Name())
		}

		name := strings.TrimSuffix(file.Name(), suffixInputFile)
		names = append(names, name)
	}

	return names, nil
}

func GoldenFiles(t *testing.T, convert ConvertFunc, roundTripConvert ConvertFunc) {
	pathOfFolder := filepath.Join("./testdata", strings.TrimPrefix(t.Name(), "Test"))
	runs, err := getInputFiles(pathOfFolder)
	if err != nil {
		t.Fatal(err)
	}
	if len(runs) == 0 {
		t.Fatalf("there were no golden files found in %q", pathOfFolder)
	}

	for _, run := range runs {
		t.Run(run, func(t *testing.T) {
			pathOfFile := filepath.Join(pathOfFolder, run+suffixInputFile)
			input, err := os.ReadFile(pathOfFile)
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("running golden file test for %q", pathOfFile)

			// - - - - - - - Golden File Test - - - - - - - //
			output, err := convert(input)
			if err != nil {
				t.Fatal(err)
			}

			g := goldie.New(t,
				goldie.WithFixtureDir(pathOfFolder),
				goldie.WithNameSuffix(suffixOutputFile),
				// Simple, ColoredDiff, ClassicDiff
				// goldie.WithDiffEngine(goldie.Simple),
				goldie.WithDiffFn(func(actual, expected string) string {
					return fmt.Sprintf("Expected: %q\nGot: %q", expected, actual)
				}),
			)
			g.Assert(t, run, []byte(output))

			// - - - - - - - Round Trip Test - - - - - - - //
			if *enableRoundTrip {
				_, err := RoundTrip(run, input, roundTripConvert)
				if err != nil {
					t.Error(err)

					// - - - //

					// TODO: clear folder to override earlier runs
					// TODO: enable writing using command line
					// err2 := res.WriteToFiles("./.tmp/roundtrip")
					// if err2 != nil {
					// 	t.Error(err2)
					// }
				}
			}
		})
	}
}
