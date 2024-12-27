package main

import (
	"os"

	"github.com/JohannesKaufmann/html-to-markdown/v2/cli/html2markdown/cmd"
)

var (
	// These are set by goreleaser:
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
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
