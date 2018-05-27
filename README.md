# html-to-markdown

![gopher stading on top of a machine that converts a box of html to blocks of markdown](/logo.png)


Convert HTML into Markdown with Go.

TODO: properly use options, list of features, document ConvertX functions

## Installation

```
go get github.com/JohannesKaufmann/html-to-markdown
```

## Usage

```go
import "github.com/JohannesKaufmann/html-to-markdown"

converter := md.NewConverter("", true, nil)

html = `<strong>Important</strong>`

markdown, err := converter.ConvertString(html)
if err != nil {
  log.Fatal(err)
}
fmt.Println("md ->", markdown)
```
If you are already using [goquery](https://github.com/PuerkitoBio/goquery) you can pass a selection to `Convert`.

```go
markdown, err := converter.Convert(selec)
```

## Adding Rules

```go
converter.AddRules(
  md.Rule{
    Filter: []string{"del", "s", "strike"},
    Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
      // You need to return a pointer to a string (md.String is just a helper function).
      // If you return nil the next function for that html element 
      // will be picked. For example you could only convert an element
      // if it has a certain class name and fallback if not.
      content = strings.TrimSpace(content)
      return md.String("~" + content + "~")
    },
  },
  // more rules
)
```

For more information have a look at the example [add_rules](/examples/add_rules/main.go).

## Using Plugins

If you want plugins (github flavored markdown like striketrough, tables, ...) you can pass it to `Use`.

```go
import "github.com/JohannesKaufmann/html-to-markdown/plugin"
  
// Use the `GitHubFlavored` plugin from the `plugin` package.
converter.Use(plugin.GitHubFlavored())
```

Or if you only want to use the `Strikethrough` plugin. You can change the character that distinguishes
the text that is crossed out by setting the first argument to a different value (for example "~~" instead of "~").
```go
converter.Use(plugin.Strikethrough(""))
```
For more information have a look at the example [github_flavored](/examples/github_flavored/main.go).

## Writing Plugins

Have a look at the [plugin folder](/plugin) for a reference implementation. The most basic one is [Strikethrough](/plugin/strikethrough.go).


## Other Methods

[Godoc](https://godoc.org/github.com/JohannesKaufmann/html-to-markdown)

### `func (c *Converter) Keep(tags ...string) *Converter`

Determines which elements are to be kept and rendered as HTML.

### `func (c *Converter) Remove(tags ...string) *Converter`

Determines which elements are to be removed altogether i.e. converted to an empty string. 


## Related Projects
- [turndown (js)](https://github.com/domchristie/turndown), a very good library written in javascript.
- [lunny/html2md](https://github.com/lunny/html2md), which is using [regex instead of goquery](https://stackoverflow.com/a/1732454). I came around a few edge case when using it (leaving some html comments, ...) so I wrote my own. 
