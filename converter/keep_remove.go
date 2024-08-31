package converter

type tagStrategy string

const (
	// - - - - - removing - - - - - //

	// StrategyRemoveNode will remove that node in the _PreRender_ phase
	// with a high priority.
	StrategyRemoveNode tagStrategy = "StrategyRemoveNode"

	// - - - - - markdown - - - - - //

	// StrategyMarkdownLeaf will keep the children of this node as markdown.
	//
	// This is the default for unknown nodes — where there
	// is no registered _Render_ handler.
	StrategyMarkdownLeaf tagStrategy = "StrategyMarkdownLeaf"

	// StrategyMarkdownBlock will keep the children of this node as markdown
	// AND will render newlines.
	//
	// This is the default for html nodes that have
	// a) no registered _Render_ handler AND
	// b) where `dom.NameIsBlockNode()` returns true.
	StrategyMarkdownBlock tagStrategy = "StrategyMarkdownBlock"

	// - - - - - html - - - - - //

	// TODO: is this needed?
	// StrategyHTMLLeaf will render the node as HTML using `html.Render()`
	// StrategyHTMLLeaf tagStrategy = "StrategyHTMLLeaf"

	// StrategyHTMLBlock will render the node as HTML using `html.Render()`
	StrategyHTMLBlock tagStrategy = "StrategyHTMLBlock"

	// - - - - - html & markdown - - - - - //

	// StrategyHTMLBlockWithMarkdown will render the node as HTML
	// and render the children as markdown.
	StrategyHTMLBlockWithMarkdown tagStrategy = "StrategyHTMLBlockWithMarkdown"
)
