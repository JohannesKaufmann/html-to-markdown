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
	// This makes `go install` and `go build` report some appropriate fallback values.
	if info, ok := debug.ReadBuildInfo(); ok {
		if version == "unknown" && info.Main.Version != "" && info.Main.Version != "(devel)" {
			version = info.Main.Version
		}
		for _, s := range info.Settings {
			switch s.Key {
			case "vcs.revision":
				if commit == "unknown" && len(s.Value) >= 7 {
					commit = s.Value[:7]
				}
			case "vcs.time":
				if date == "unknown" && s.Value != "" {
					date = s.Value // e.g. "2026-05-07T19:14:37Z"
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
