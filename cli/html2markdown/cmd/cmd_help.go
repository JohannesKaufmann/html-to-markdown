package cmd

import (
	"flag"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/template"
)

var usageTemplate = `
# html2markdown - convert html to markdown [version {{ .Version }}]

Convert HTML to Markdown. Even works with entire websites!

## Basics

By default the "Commonmark" Plugin will be enabled. You can customize the options,
for example changing the appearance of bold with --opt-strong-delimiter="__"

Other Plugins can also be enabled. For example "GitHub Flavored Markdown" (GFM)
extends Commonmark with more features.

## Relative / Absolute Links

Use --domain="https://example.com" to convert *relative* links to *absolute* links.
The same also works for images.

## Escaping

Some characters have a special meaning in markdown. The library escapes these — if necessary.
See the documentation for more info.

## Security

Once you convert this markdown *back* to HTML you need to be careful of malicious content. 
Use a HTML sanitizer before displaying the HTML in the browser!


## Examples

    echo "<strong>important</strong>" | html2markdown

    curl --no-progress-meter http://example.com | html2markdown


    html2markdown --input file.html --output file.md

    html2markdown --input "src/*.html" --output "dist/"


## Flags

    -v, --version
        show the version of html2markdown and exit

    --help

    --input PATH
        Input file, directory, or glob pattern (instead of stdin)

    --output PATH
        Output file or directory (instead of stdout)

    --output-overwrite
        Replace existing files

    If --input is a directory or glob pattern, --output must be a directory.


{{ range .Flags }}
    --{{ .Name }}{{ with .Usage }}
{{ . | indent 8 }}{{ end }}
{{ end }}


For more information visit the documentation:
https://github.com/Johanneskaufmann/html-to-markdown

`

var templateFuncs = template.FuncMap{
	"indent": func(spaces int, v string) string {
		pad := strings.Repeat(" ", spaces)
		return pad + strings.Replace(v, "\n", "\n"+pad, -1)
	},
}

func tmpl(w io.Writer, text string, data interface{}) error {
	t := template.New("usage")
	t.Funcs(templateFuncs)

	_, err := t.Parse(text)
	if err != nil {
		return err
	}
	return t.Execute(w, data)
}

func (cli *CLI) initUsageText() error {
	var flags []*flag.Flag
	cli.flags.VisitAll(func(f *flag.Flag) {
		if f.Name == "v" || f.Name == "version" || f.Name == "input" || f.Name == "output" || f.Name == "output-overwrite" {
			// We manually mention these in the usage
			return
		}
		flags = append(flags, f)
	})
	sort.Slice(flags, func(i, j int) bool {
		return flags[i].Name < flags[j].Name
	})

	data := map[string]any{
		"Version": cli.Release.Version,
		"Flags":   flags,
	}
	err := tmpl(&cli.usageText, usageTemplate, data)
	if err != nil {
		return err
	}

	return nil
}

func (cli CLI) printUsage() {
	fmt.Fprint(cli.Stdout, cli.usageText.String())
}
