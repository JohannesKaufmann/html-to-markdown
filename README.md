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

---

---

## CLI - Using it on the command line

Using the Golang library provides the most customization, while the CLI is the simplest way to get started.

### Installation

Download the pre-compiled binaries from the [releases page](https://github.com/JohannesKaufmann/html-to-markdown/releases) and copy them to the desired location.

```bash
html2markdown --version
```

> [!NOTE]  
> Make sure that `--version` prints `2.X.X` as there is a different CLI for V2 of the converter.

## Usage

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
