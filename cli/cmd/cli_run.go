package cmd

func Run(
	stdin ReadWriterWithStat,
	stdout ReadWriterWithStat,
	stderr ReadWriterWithStat,

	osArgs []string,

	release Release,
) {

	cli := CLI{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,

		OsArgs: osArgs,

		Release: release,
	}

	// - - - - - init - - - - - //
	if err := cli.Init(); err != nil {
		panic(err)
	}

	// - - - - - exec - - - - - //
	cli.Execute()
}
