package feed

import (
	"fmt"
	gorFeed "github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
	"io/ioutil"
	"net/http"
	"time"
)

// Feed represents our base unit of work
type Feed struct {
	fp              *gofeed.Parser
	URL             string
	RawData         *http.Response
	OriginalHeaders http.Header
	OriginalContent []byte
	ParsedContent   *gofeed.Feed
	FilteredItems   []*gofeed.Item
}

// New creates a new feed object
func New(feedURL string) (*Feed, error) {
	return &Feed{
		URL: feedURL,
		fp:  gofeed.NewParser(),
	}, nil
}

// Download pulls down a specified feed
func (f *Feed) Download() error {
	response, err := http.Get(f.URL)
	if err != nil {
		return err
	}
	f.RawData = response
	f.OriginalHeaders = response.Header
	defer response.Body.Close()
	f.OriginalContent, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return nil
}

// Parse a feed into a gofeed.Feed data structure
func (f *Feed) Parse() error {
	fp := gofeed.NewParser()
	content, err := fp.ParseString(string(f.OriginalContent))
	if err != nil {
		return err
	}
	f.ParsedContent = content
	return nil
}

// Filter a feed using a rulefile
func (f *Feed) Filter() error {
	if f.ParsedContent == nil {
		return fmt.Errorf("Feed not parsed yet")
	}

	for _, item := range f.ParsedContent.Items {
		f.FilteredItems = append(f.FilteredItems, item)
	}
	return nil
}

// ToXML writes the filtered feed back into XML
func (f *Feed) ToXML() (string, error) {
	now := time.Now()

	if f.ParsedContent == nil {
		return "", fmt.Errorf("Feed not parsed yet")
	}

	feed := &gorFeed.Feed{
		Title:       f.ParsedContent.Title,
		Link:        &gorFeed.Link{Href: "https://localhost"},
		Description: "RSS filtering",
		Author:      &gorFeed.Author{Name: "John Doe", Email: "jdoe@example.org"},
		Created:     now,
	}

	for _, item := range f.FilteredItems {
		entry := &gorFeed.Item{
			Id:      item.GUID,
			Title:   item.Title,
			Link:    &gorFeed.Link{Href: item.Link},
			Created: *item.UpdatedParsed,
			Content: item.Content,
		}
		if item.Author != nil {
			entry.Author = &gorFeed.Author{Name: item.Author.Name, Email: item.Author.Email}
		}
		if len(item.Description) > 0 {
			entry.Description = item.Description
		}

		feed.Add(entry)
	}

	feedXML, err := feed.ToAtom()
	if err != nil {
		return "", err
	}

	return feedXML, nil
}
