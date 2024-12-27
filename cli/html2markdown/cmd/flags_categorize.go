package cmd

import (
	"flag"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/agnivade/levenshtein"
)

const flagProvidedButNotDefinedErr = "flag provided but not defined: -"

func formatFlag(name string) string {
	if len(name) == 1 {
		return "-" + name
	}
	return "--" + name
}
func (cli *CLI) getAlternativeFlag(unknownFlag string) string {
	var closestDistance int = 10000
	var closestFlag string

	cli.flags.VisitAll(func(f *flag.Flag) {

		distance := levenshtein.ComputeDistance(f.Name, unknownFlag)

		if distance < closestDistance {
			closestDistance = distance
			closestFlag = f.Name
		}
	})

	if closestDistance >= utf8.RuneCountInString(unknownFlag) {
		return ""
	}
	if closestDistance > 4 {
		return ""
	}
	return closestFlag
}
func (cli *CLI) categorizeFlagError(err error) error {
	if err == nil {
		return nil
	}

	message := err.Error()

	if strings.HasPrefix(message, flagProvidedButNotDefinedErr) {
		flagName := strings.TrimPrefix(message, flagProvidedButNotDefinedErr)

		err := fmt.Errorf("unknown flag: %s", formatFlag(flagName))

		alternative := cli.getAlternativeFlag(flagName)
		if alternative == "" {
			return NewCLIError(err)
		}

		return NewCLIError(
			err,
			Paragraph(fmt.Sprintf("Did you mean %s instead?", formatFlag(alternative))),
		)
	}

	return err
}
