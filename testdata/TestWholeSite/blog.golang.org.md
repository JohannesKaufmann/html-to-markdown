Godoc: documenting Go code - The Go Blog[The Go Programming Language](//golang.org/)[Go](//golang.org/)[▽](#)[Documents](//golang.org/doc/)[Packages](//golang.org/pkg/)[The Project](//golang.org/project/)[Help](//golang.org/help/)[Blog](/)submit search

#### Next article

[Introducing Gofix](/introducing-gofix)

#### Previous article

[Gobs of data](/gobs-of-data)

#### Links

- [golang.org](//golang.org/)
- [Install Go](//golang.org/doc/install.html)
- [A Tour of Go](//tour.golang.org/)
- [Go Documentation](//golang.org/doc/)
- [Go Mailing List](//groups.google.com/group/golang-nuts)
- [Go on Google+](//plus.google.com/101406623878176903605)
- [Go+ Community](//plus.google.com/communities/114112804251407510571)
- [Go on Twitter](//twitter.com/golang)

[Blog index](/index)

# [The Go Blog](/)

### [Godoc: documenting Go code](/godoc-documenting-go-code)

31 March 2011

 The Go project takes documentation seriously. Documentation is a huge part of making software accessible and maintainable.
 Of course it must be well-written and accurate, but it also must be easy to write and to maintain. Ideally, it
 should be coupled to the code itself so the documentation evolves along with the code. The easier it is for programmers
 to produce good documentation, the better for everyone.
 

 To that end, we have developed the
 [godoc](https://golang.org/cmd/godoc/) documentation tool. This article describes godoc's approach to documentation, and explains how
 you can use our conventions and tools to write good documentation for your own projects.
 

 Godoc parses Go source code - including comments - and produces documentation as HTML or plain text. The end result is documentation
 tightly coupled with the code it documents. For example, through godoc's web interface you can navigate from
 a function's
 [documentation](https://golang.org/pkg/strings/#HasPrefix) to its
 [implementation](https://golang.org/src/pkg/strings/strings.go#L493) with one click.
 

 Godoc is conceptually related to Python's
 [Docstring](http://www.python.org/dev/peps/pep-0257/) and Java's
 [Javadoc](http://www.oracle.com/technetwork/java/javase/documentation/index-jsp-135444.html), but its design is simpler. The comments read by godoc are not language constructs (as with Docstring)
 nor must they have their own machine-readable syntax (as with Javadoc). Godoc comments are just good comments,
 the sort you would want to read even if godoc didn't exist.
 

 The convention is simple: to document a type, variable, constant, function, or even a package, write a regular comment directly
 preceding its declaration, with no intervening blank line. Godoc will then present that comment as text alongside
 the item it documents. For example, this is the documentation for the
 `fmt` package's
 [`Fprint`](https://golang.org/pkg/fmt/#Fprint) function:
 

```
// Fprint formats using the default formats for its operands and writes to w.
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func Fprint(w io.Writer, a ...interface{}) (n int, err error) {
```

 Notice this comment is a complete sentence that begins with the name of the element it describes. This important convention
 allows us to generate documentation in a variety of formats, from plain text to HTML to UNIX man pages, and makes
 it read better when tools truncate it for brevity, such as when they extract the first line or sentence.
 

 Comments on package declarations should provide general package documentation. These comments can be short, like the
 [`sort`](https://golang.org/pkg/sort/) package's brief description:
 

```
// Package sort provides primitives for sorting slices and user-defined
// collections.
package sort
```

 They can also be detailed like the
 [gob package](https://golang.org/pkg/encoding/gob/)'s overview. That package uses another convention for packages that need large amounts of
 introductory documentation: the package comment is placed in its own file,
 [doc.go](https://golang.org/src/pkg/encoding/gob/doc.go), which contains only those comments and a package clause.
 

 When writing package comments of any size, keep in mind that their first sentence will appear in godoc's
 [package list](https://golang.org/pkg/).
 

 Comments that are not adjacent to a top-level declaration are omitted from godoc's output, with one notable exception.
 Top-level comments that begin with the word
 `"BUG(who)”` are recognized as known bugs, and included in the "Bugs” section of the package documentation. The "who”
 part should be the user name of someone who could provide more information. For example, this is a known issue
 from the
 [bytes package](https://golang.org/pkg/bytes/#pkg-note-BUG):
 

```
// BUG(r): The rule Title uses for word boundaries does not handle Unicode punctuation properly.
```

 Sometimes a struct field, function, type, or even a whole package becomes redundant or unnecessary, but must be kept for
 compatibility with existing programs. To signal that an identifier should not be used, add a paragraph to its
 doc comment that begins with "Deprecated:" followed by some information about the deprecation. There
 are a few examples
 [in the standard library](https://golang.org/search?q=Deprecated:).
 

 There are a few formatting rules that Godoc uses when converting comments to HTML:
 

- Subsequent lines of text are considered part of the same paragraph; you must leave a blank line to separate paragraphs.

- Pre-formatted text must be indented relative to the surrounding comment text (see gob's
     [doc.go](https://golang.org/src/pkg/encoding/gob/doc.go) for an example).

- URLs will be converted to HTML links; no special markup is necessary.

 Note that none of these rules requires you to do anything out of the ordinary.
 

 In fact, the best thing about godoc's minimal approach is how easy it is to use. As a result, a lot of Go code, including
 all of the standard library, already follows the conventions.
 

 Your own code can present good documentation just by having comments as described above. Any Go packages installed inside
 `$GOROOT/src/pkg` and any
 `GOPATH` work spaces will already be accessible via godoc's command-line and HTTP interfaces, and you can specify
 additional paths for indexing via the
 `-path` flag or just by running
 `"godoc ."` in the source directory. See the
 [godoc documentation](https://golang.org/cmd/godoc/) for more details.
 

By Andrew Gerrand

## Related articles

- [HTTP/2 Server Push](/h2push)
- [Introducing HTTP Tracing](/http-tracing)
- [Testable Examples in Go](/examples)
- [Generating code](/generate)
- [Introducing the Go Race Detector](/race-detector)
- [Go maps in action](/go-maps-in-action)
- [go fmt your code](/go-fmt-your-code)
- [Organizing Go code](/organizing-go-code)
- [Debugging Go programs with the GNU Debugger](/debugging-go-programs-with-gnu-debugger)
- [The Go image/draw package](/go-imagedraw-package)
- [The Go image package](/go-image-package)
- [The Laws of Reflection](/laws-of-reflection)
- [Error handling and Go](/error-handling-and-go)
- ["First Class Functions in Go"](/first-class-functions-in-go-and-new-go)
- [Profiling Go Programs](/profiling-go-programs)
- [A GIF decoder: an exercise in Go interfaces](/gif-decoder-exercise-in-go-interfaces)
- [Introducing Gofix](/introducing-gofix)
- [Gobs of data](/gobs-of-data)
- [C? Go? Cgo!](/c-go-cgo)
- [JSON and Go](/json-and-go)
- [Go Slices: usage and internals](/go-slices-usage-and-internals)
- [Go Concurrency Patterns: Timing out, moving on](/go-concurrency-patterns-timing-out-and)
- [Defer, Panic, and Recover](/defer-panic-and-recover)
- [Share Memory By Communicating](/share-memory-by-communicating)
- [JSON-RPC: a tale of interfaces](/json-rpc-tale-of-interfaces)

 Except as
 [noted](https://developers.google.com/site-policies#restrictions), the content of this page is licensed under the Creative Commons Attribution 3.0 License,
  and code is licensed under a
 [BSD license](//golang.org/LICENSE).
 [Terms of Service](//golang.org/doc/tos.html) |
 [Privacy Policy](//www.google.com/intl/en/policies/privacy/) |
 [View the source code](https://go.googlesource.com/blog/)