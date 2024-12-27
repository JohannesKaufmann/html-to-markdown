package cmd

import (
	"bytes"
	"fmt"

	"github.com/JohannesKaufmann/dom"
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/strikethrough"
	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

func overrideValidationError(e *commonmark.ValidateConfigError) error {

	// TODO: Maybe OptionFunc should already validate and return an error?
	//       Then it would be easier to override the Key since we have once
	//       place to assemble the []OptionFunc and directly treat the errors...
	//
	// We would basically invoke it ourselves:
	//    err := commonmark.WithStrongDelimiter(cli.config.strongDelimiter)(conv)

	switch e.Key {
	case "StrongDelimiter":
		e.Key = "opt-strong-delimiter"
	}

	e.KeyWithValue = fmt.Sprintf("--%s=%q", e.Key, e.Value)
	return e
}

func (cli *CLI) includeNodesFromDoc(doc *html.Node) (*html.Node, error) {
	if len(cli.config.includeSelector) == 0 {
		return doc, nil
	}
	nodes := cascadia.QueryAll(doc, cli.config.includeSelector)

	root := &html.Node{}
	for _, n := range nodes {
		dom.RemoveNode(n)
		root.AppendChild(n)
	}

	return root, nil
}
func (cli *CLI) excludeNodesFromDoc(doc *html.Node) error {
	if len(cli.config.excludeSelector) == 0 {
		return nil
	}

	var finder func(node *html.Node)
	finder = func(node *html.Node) {
		if cli.config.excludeSelector.Match(node) {
			dom.RemoveNode(node)
			return
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			// Because we are sometimes removing a node, this causes problems
			// with the for loop. Using `defer` is a cool trick!
			// https://gist.github.com/loopthrough/17da0f416054401fec355d338727c46e
			defer finder(child)
		}
	}
	finder(doc)

	return nil
}
func (cli *CLI) parseInputWithSelectors(input []byte) (*html.Node, error) {
	r := bytes.NewReader(input)

	doc, err := html.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("error while parsing html: %w", err)
	}

	doc, err = cli.includeNodesFromDoc(doc)
	if err != nil {
		return nil, err
	}

	err = cli.excludeNodesFromDoc(doc)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (cli *CLI) convert(input []byte) ([]error, error) {
	conv := converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),
			commonmark.NewCommonmarkPlugin(
				commonmark.WithStrongDelimiter(cli.config.strongDelimiter),
			),
		),
	)
	if cli.config.enablePluginStrikethrough {
		// TODO: while this works, this does not add the `Name` to the internal list
		strikethrough.NewStrikethroughPlugin().Init(conv)
	}

	doc, err := cli.parseInputWithSelectors(input)
	if err != nil {
		return nil, err
	}

	markdown, err := conv.ConvertNode(doc, converter.WithDomain(cli.config.domain))
	if err != nil {
		e, ok := err.(*commonmark.ValidateConfigError)
		if ok {
			return nil, overrideValidationError(e)
		}

		return nil, err
	}

	fmt.Fprintln(cli.Stdout, string(markdown))
	return nil, nil
}
