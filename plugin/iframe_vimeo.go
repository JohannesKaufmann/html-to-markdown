package plugin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

// Timeout for the http client
var Timeout = time.Second * 10
var netClient = &http.Client{
	Timeout: Timeout,
}

type vimeoVideo struct {
	Type                       string `json:"type"`
	Version                    string `json:"version"`
	ProviderName               string `json:"provider_name"`
	ProviderURL                string `json:"provider_url"`
	Title                      string `json:"title"`
	AuthorName                 string `json:"author_name"`
	AuthorURL                  string `json:"author_url"`
	IsPlus                     string `json:"is_plus"`
	AccountType                string `json:"account_type"`
	HTML                       string `json:"html"`
	Width                      int    `json:"width"`
	Height                     int    `json:"height"`
	Duration                   int    `json:"duration"`
	Description                string `json:"description"`
	ThumbnailURL               string `json:"thumbnail_url"`
	ThumbnailWidth             int    `json:"thumbnail_width"`
	ThumbnailHeight            int    `json:"thumbnail_height"`
	ThumbnailURLWithPlayButton string `json:"thumbnail_url_with_play_button"`
	UploadDate                 string `json:"upload_date"`
	VideoID                    int    `json:"video_id"`
	URI                        string `json:"uri"`
}

var vimeoID = regexp.MustCompile(`video\/(\d*)`)

type vimeoVariation int

// Configure how the Vimeo Plugin should display the video in markdown.
const (
	VimeoOnlyThumbnail vimeoVariation = iota
	VimeoWithTitle
	VimeoWithDescription
)

// VimeoEmbed registers a rule (for iframes) and
// returns a markdown compatible representation (link to video, ...).
func VimeoEmbed(variation vimeoVariation) md.Plugin {
	return func(c *md.Converter) []md.Rule {
		getVimeoData := func(id string) (*vimeoVideo, error) {
			u := fmt.Sprintf("http://vimeo.com/api/oembed.json?url=https://vimeo.com/%s", id)

			resp, err := netClient.Get(u)
			if err != nil {
				return nil, err
			}

			defer resp.Body.Close()

			var res vimeoVideo
			err = json.NewDecoder(resp.Body).Decode(&res)
			if err != nil {
				return nil, err
			}
			return &res, nil
		}
		cleanDescription := func(html string) (string, error) {
			text, err := c.ConvertString(html)
			if err != nil {
				return "", err
			}

			text = strings.Replace(text, "\n", " ", -1)
			text = strings.Replace(text, "\t", " ", -1)
			before := utf8.RuneCountInString(text)
			text = summary(text, 70)
			after := utf8.RuneCountInString(text)
			if after != before {
				text += "..."
			}
			return text, nil
		}

		return []md.Rule{
			{
				Filter: []string{"iframe"},
				Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
					src := selec.AttrOr("src", "")
					if !strings.Contains(src, "vimeo.com") {
						return nil
					}
					parts := vimeoID.FindStringSubmatch(src)
					if len(parts) != 2 {
						return nil
					}
					id := parts[1]

					video, err := getVimeoData(id)
					if err != nil {
						panic(err)
					}

					// desc, err := cleanDescription(video.Description)
					// if err != nil {
					// 	panic(err)
					// }

					// [![Little red ridning hood](http://i.imgur.com/7YTMFQp.png)](https://vimeo.com/3514904 "Little red riding hood - Click to Watch!")
					// text := fmt.Sprintf("[![%s](%s) ](%s)", desc, video.ThumbnailURLWithPlayButton, "https://vimeo.com/"+video.URI)
					text := fmt.Sprintf(`[![](%s)](https://vimeo.com/%d)`, video.ThumbnailURLWithPlayButton, video.VideoID)

					switch variation {
					case VimeoOnlyThumbnail:
						// do nothing
					case VimeoWithTitle:
						duration := time.Duration(video.Duration) * time.Second
						text += fmt.Sprintf("\n\n'%s' by ['%s'](%s) (%s)", video.Title, video.AuthorName, video.AuthorURL, duration.String())
					case VimeoWithDescription:
						desc, err := cleanDescription(video.Description)
						if err != nil {
							panic(err)
						}
						text += "\n\n" + desc
					}

					return &text
				},
			},
		}
	}
}

// truncate
func summary(text string, limit int) string {
	result := text
	chars := 0
	for i := range text {
		if chars >= limit {
			result = text[:i]
			break
		}
		chars++
	}
	return result
}
