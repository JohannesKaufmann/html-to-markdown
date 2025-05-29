package cmd

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/JohannesKaufmann/dom"
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/strikethrough"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/table"
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

func (cli *CLI) convert(input []byte) ([]byte, error) {
	conv := converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),
			commonmark.NewCommonmarkPlugin(
				commonmark.WithStrongDelimiter(cli.config.strongDelimiter),
			),
		),
	)
	if cli.config.enablePluginStrikethrough {
		conv.Register.Plugin(strikethrough.NewStrikethroughPlugin())
	}

	if cli.config.enablePluginTable {
		conv.Register.Plugin(
			table.NewTablePlugin(
				table.WithSkipEmptyRows(cli.config.tableSkipEmptyRows),
				table.WithHeaderPromotion(cli.config.tableHeaderPromotion),
				table.WithSpanCellBehavior(table.SpanCellBehavior(cli.config.tableSpanCellBehavior)),
				table.WithPresentationTables(cli.config.tablePresentationTables),
				table.WithNewlineBehavior(table.NewlineBehavior(cli.config.tableNewlineBehavior)),
				table.WithPadColumns(cli.config.tablePadColumns),
			),
		)
	}

	doc, err := cli.parseInputWithSelectors(input)
	if err != nil {
		return nil, err
	}

	markdown, err := conv.ConvertNode(doc, converter.WithDomain(cli.config.domain))
	if err != nil {

		var validationErr *commonmark.ValidateConfigError
		if errors.As(err, &validationErr) {
			return nil, overrideValidationError(validationErr)
		}

		return nil, err
	}

	return markdown, nil
}
