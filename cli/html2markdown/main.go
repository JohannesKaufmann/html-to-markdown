package main

import (
	"os"
	"runtime/debug"

	"github.com/JohannesKaufmann/html-to-markdown/v2/cli/html2markdown/cmd"
)

var (
	// These are set by goreleaser:
	version = "unknown"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	// Fall back to build info when ldflags are not set by goreleaser.
	// This makes `go install` report the correct module version and
	// `go build` report the VCS commit hash.
	if info, ok := debug.ReadBuildInfo(); ok {
		if version == "unknown" && info.Main.Version != "" && info.Main.Version != "(devel)" {
			version = info.Main.Version
		}
		if commit == "unknown" {
			for _, s := range info.Settings {
				if s.Key == "vcs.revision" && len(s.Value) >= 7 {
					commit = s.Value[:7]
					break
				}
			}
		}
	}

	release := cmd.Release{
		Version: version,
		Commit:  commit,
		Date:    date,
	}

	cmd.Run(
		os.Stdin,
		os.Stdout,
		os.Stderr,
		os.Args,
		release,
	)
}
