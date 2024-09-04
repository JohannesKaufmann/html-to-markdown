# html-to-markdown

> [!WARNING]
> This is an **early experimental version** of the library.
>
> We encourage testing and bug reporting. However, please note:
>
> - Not production-ready
>   - Default options are well-tested, but custom configurations have limited coverage
> - Functionality is currently restricted
>   - Focus is on stabilization and core features
> - No compatibility guarantee
>   - Only use `htmltomarkdown.ConvertString()` and `htmltomarkdown.ConvertNode()` from the root package. They are _unlikely_ to change.
>   - Other functions and nested packages are _very like_ to change.

---

## Golang Library

### Installation

```bash
go get -u github.com/JohannesKaufmann/html-to-markdown/v2
```

_Or if you want a specific commit add the suffix `/v2@commithash`_

### Usage

[![Go V2 Reference](https://pkg.go.dev/badge/github.com/JohannesKaufmann/html-to-markdown/v2.svg)](https://pkg.go.dev/github.com/JohannesKaufmann/html-to-markdown/v2)

```go
package main

import (
	"fmt"
	"log"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
)

func main() {
	input := `<strong>Bold Text</strong>`

	markdown, err := htmltomarkdown.ConvertString(input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(markdown)
	// Output: **Bold Text**
}
```

- 🧑‍💻 [Example code, basics](/examples/basics/main.go)

The function `htmltomarkdown.ConvertString()` is just a small wrapper around `converter.NewConverter()` and `commonmark.NewCommonmarkPlugin()`. If you want more control, use the following:

```go
package main

import (
	"fmt"
	"log"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
)

func main() {
	input := `<strong>Bold Text</strong>`

	conv := converter.NewConverter(
		converter.WithPlugins(
			commonmark.NewCommonmarkPlugin(
				commonmark.WithStrongDelimiter("__"),
				// ...additional configurations for the plugin
			),
		),
	)

	markdown, err := conv.ConvertString(input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(markdown)
	// Output: __Bold Text__
}
```

- 🧑‍💻 [Example code, options](/examples/options/main.go)

> [!NOTE]  
> If you use `NewConverter` directly make sure to also **register the commonmark plugin**.

### Plugins

TODO: info about plugins

---

---

## CLI - Using it on the command line

Using the Golang library provides the most customization, while the CLI is the simplest way to get started.

### Installation

#### Homebrew Tap

```bash
brew install JohannesKaufmann/tap/html2markdown
```

#### Manually

Download the pre-compiled binaries from the [releases page](https://github.com/JohannesKaufmann/html-to-markdown/releases) and copy them to the desired location.

### Version

```bash
html2markdown --version
```

> [!NOTE]  
> Make sure that `--version` prints `2.X.X` as there is a different CLI for V2 of the converter.

### Usage

```bash
$ echo "<strong>important</strong>" | html2markdown

**important**
```

```text
$ curl --no-progress-meter http://example.com | html2markdown

# Example Domain

This domain is for use in illustrative examples in documents. You may use this domain in literature without prior coordination or asking for permission.

[More information...](https://www.iana.org/domains/example)
```

_(The cli does not support every option yet. Over time more customization will be added)_

---

---

## FAQ

### Extending with Plugins

- Need your own logic? Write your own code and then **register** it.
- Don't like the **defaults** that the library uses? You can use `PriorityEarly` to run you logic _earlier_ than others.
- If you believe that you logic could also benefit others, you can package it up into a **plugin**.

🗒️ [WRITING_PLUGINS.md](/WRITING_PLUGINS.md)

### Bugs

You found a bug?

[Open an issue](https://github.com/JohannesKaufmann/html-to-markdown/issues/new/choose) with the HTML snippet that does not produce the expected results. Please, please, plase _submit the HTML snippet_ that caused the problem. Otherwise it is very difficult to reproduce and fix...

### Security

This library produces markdown that is readable and can be changed by humans.

Once you convert this markdown back to HTML (e.g. using [goldmark](https://github.com/yuin/goldmark) or [blackfriday](https://github.com/russross/blackfriday)) you need to be careful of malicious content.

This library does NOT sanitize untrusted content. Use an HTML sanitizer such as [bluemonday](https://github.com/microcosm-cc/bluemonday) before displaying the HTML in the browser.

🗒️ [SECURITY.md](/SECURITY.md) if you find a security vulnerability

### Goroutines

You can use the `Converter` from (multiple) goroutines. Internally a mutex is used & there is a test to verify that behaviour.

### Escaping & Backslash

Some characters have a special meaning in markdown (e.g. "\*" for emphasis). The backslash `\` character is used to "escape" those characters. That is perfectly safe and won't be displayed in the final render.

🗒️ [ESCAPING.md](/ESCAPING.md)

### Contributing

You want to contribute? Thats great to hear! There are many ways to help:

Helping to answer questions, triaging issues, writing documentation, writing code, ...

If you want to make a code change: Please first discuss the change you wish to make, by opening an issue. I'm also happy to guide you to where a change is most likely needed. There are also extensive tests (see below) so you can freely experiment 🧑‍🔬

_Note: The outside API should not change because of backwards compatibility..._

### Testing

You don't have to be afraid of breaking the converter, since there are many "Golden File" tests:

Add your problematic HTML snippet to one of the `.in.html` files in the `testdata` folders. Then run `go test -update` and have a look at which `.out.md` files changed in GIT.

You can now change the internal logic and inspect what impact your change has by running `go test -update` again.

_Note: Before submitting your change as a PR, make sure that you run those tests and check the files into GIT..._

### License

Unless otherwise specified, the project is licensed under the terms of the MIT license.

🗒️ [LICENSE](/LICENSE)
