package tester

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	// "github.com/darmiel/gohtml"
	"github.com/yuin/goldmark"
	goldmarkHtml "github.com/yuin/goldmark/renderer/html"
)

type ConvertFunc func(html []byte) (markdown []byte, err error)

var goldmarkConverter = goldmark.New(
	goldmark.WithRendererOptions(
		// Also render "javascript:" links
		goldmarkHtml.WithUnsafe(),
	),
)

type Result struct {
	Identifier string

	FirstDuration  time.Duration
	SecondDuration time.Duration

	OriginalHtml     []byte
	FirstMarkdown    []byte
	IntermediateHtml []byte
	SecondMarkdown   []byte
}

func (r Result) GetStatus() string {
	if bytes.Equal(r.FirstMarkdown, r.SecondMarkdown) {
		return fmt.Sprintf("%s: ✅ in %s & %s", r.Identifier, r.FirstDuration, r.SecondDuration)
	}

	return fmt.Sprintf("%s: ❌ in %s & %s", r.Identifier, r.FirstDuration, r.SecondDuration)
}
func (r Result) PrintStatus() {
	var Reset = "\033[0m"
	var Yellow = "\033[33m"
	fmt.Printf("%s[Round Trip Test] %s \n%s", Yellow, r.GetStatus(), Reset)

}
func (r Result) WriteToFiles(folderpath string) error {
	err := os.MkdirAll(filepath.Join(folderpath, r.Identifier), os.ModePerm)
	if err != nil {
		return fmt.Errorf("error while creating folder %q: %w", filepath.Join(folderpath, r.Identifier), err)
	}

	err = os.WriteFile(filepath.Join(folderpath, r.Identifier, "01.html"), r.OriginalHtml, 0644)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(folderpath, r.Identifier, "02.md"), r.FirstMarkdown, 0644)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(folderpath, r.Identifier, "03.html"), r.IntermediateHtml, 0644)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(folderpath, r.Identifier, "04.md"), r.SecondMarkdown, 0644)
	if err != nil {
		return err
	}

	// originalHtmlPretty := gohtml.Format(string(r.OriginalHtml), true)
	// err = ioutil.WriteFile(filepath.Join(folderpath, r.Identifier, "01_pretty.html"), []byte(originalHtmlPretty), 0644)
	// if err != nil {
	// 	return err
	// }

	// intermediateHtmlPretty := gohtml.Format(string(r.IntermediateHtml), true)
	// err = ioutil.WriteFile(filepath.Join(folderpath, r.Identifier, "03_pretty.html"), []byte(intermediateHtmlPretty), 0644)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func RoundTrip(identifier string, originalHtml []byte, convert ConvertFunc) (*Result, error) {
	var err error
	res := &Result{
		Identifier:   identifier,
		OriginalHtml: originalHtml,
	}

	firstStart := time.Now()
	res.FirstMarkdown, err = convert(originalHtml)
	res.FirstDuration = time.Since(firstStart)
	if err != nil {
		return res, fmt.Errorf("error in the first convert round: %w", err)
	}

	var buf bytes.Buffer
	err = goldmarkConverter.Convert(res.FirstMarkdown, &buf)
	if err != nil {
		return res, fmt.Errorf("error with goldmark: %w", err)
	}
	res.IntermediateHtml = buf.Bytes()

	secondStart := time.Now()
	res.SecondMarkdown, err = convert(res.IntermediateHtml)
	res.SecondDuration = time.Since(secondStart)
	if err != nil {
		return res, fmt.Errorf("error in the second convert round: %w", err)
	}

	if bytes.Equal(res.FirstMarkdown, res.SecondMarkdown) {
		// Hurray, the converter produced exactly the same result. Well done!!!
		return res, nil
	}

	return res, fmt.Errorf("difference between the first and second markdown round")
}
