# html-to-markdown
Convert HTML into Markdown with Go.


## Installation

```
go get github.com/JohannesKaufmann/html-to-markdown
```

## Usage

```go
converter := md.NewConverter("www.google.com", true, nil)

html = `<strong>Important</strong>`

markdown, err := converter.ConvertString(html)
if err != nil {
  log.Fatal(err)
}
fmt.Println("md ->", markdown)
```
If you are already using [goquery](https://github.com/PuerkitoBio/goquery) you can pass a selection to `Convert`.


## Options

Options can be passed to `NewConverter`.

TODO: Table


## Methods

### `func NewConverter(domain string, enableCommonmark bool, options *Options) *Converter`
- `domain` is used for links and images to convert relative urls ("/image.png") to absolute urls.
    If you want more controll use TODO: OnImage/...
- [CommonMark](http://commonmark.org/) is the default set of rules. Set `enableCommonmark` to false
    if you want to customize everything using `AddRules` and DONT want to fallback to default rules.


### `func (c *Converter) AddRules(rules ...Rule) *Converter`
```go
converter.AddRules(
  md.Rule{
    Filter: []string{"del", "s", "strike"},
    Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
      // You need to return a pointer to a string (md.String is just a helper function).
      // If you return nil the next function for that html element 
      // will be picked. For example you could only convert an element
      // if it has a certain class name and fallback if not.
      return md.String("~" + content + "~")
    },
  },
  // more rules
)
```

If you want plugins (github flavored markdown like striketrough, tables, ...) you can pass it to `AddRules`.
```go
// gfm: github flavored markdown
converter.AddRules(plugin.GFM...)

// OR

converter.AddRules(plugin.Strikethrough...)
converter.AddRules(plugin.TaskListItems...)
converter.AddRules(plugin.Table...)
```

### `func (c *Converter) Keep(tags ...string) *Converter`

Determines which elements are to be kept and rendered as HTML.

### `func (c *Converter) Remove(tags ...string) *Converter`

Determines which elements are to be removed altogether i.e. converted to an empty string. 

## Related Projects
- [turndown (js)](https://github.com/domchristie/turndown), a very good library written in javascript.
- [lunny/html2md](https://github.com/lunny/html2md), which is using [regex instead of goquery](https://stackoverflow.com/a/1732454). I came around a few edge case when using it (leaving some html comments, ...) so I wrote my own. 
