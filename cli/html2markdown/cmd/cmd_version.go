package cmd

import "fmt"

func (cli CLI) printVersion() {
	fmt.Fprintf(cli.Stdout, "%s\n\n", projectBinary)

	fmt.Fprintf(cli.Stdout, "GitVersion:  %s\n", cli.Release.Version)
	fmt.Fprintf(cli.Stdout, "GitCommit:   %s\n", cli.Release.Commit)
	fmt.Fprintf(cli.Stdout, "BuildDate:   %s\n", cli.Release.Date)
}
