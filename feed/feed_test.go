package feed

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

const ()

func testFeed(t *testing.T) *Feed {
	feedXML, err := ioutil.ReadFile("../fixtures/feed.xml")
	if err != nil {
		t.Errorf("Could not read fixture")
	}

	headerJSON, err := ioutil.ReadFile("../fixtures/headers.json")
	headers := make(http.Header)
	err = json.Unmarshal(headerJSON, &headers)
	if err != nil {
		t.Errorf("Could not read header fixture: %s", err)
	}

	feed, err := New("https://localhost")
	feed.OriginalContent = feedXML
	feed.OriginalHeaders = headers

	return feed
}

func TestNew(t *testing.T) {
	_, err := New("https://localhost")
	if err != nil {
		t.Errorf("Could not create feed object")
	}
}

func TestDownload(t *testing.T) {
	feed, _ := New("https://miniflux.app/feed.xml")

	err := feed.Download()

	if err != nil || feed.RawData == nil {
		t.Errorf("Unable to download a feed")
	}

}

func TestParse(t *testing.T) {
	feed := testFeed(t)
	err := feed.Parse()

	if err != nil || feed.ParsedContent == nil {
		t.Errorf("Unable to parse a feed")
	}
}

func TestFilter(t *testing.T) {
	feed := testFeed(t)
	feed.Parse()
	err := feed.Filter()

	if err != nil || feed.FilteredItems == nil {
		t.Errorf("Unable to filter a feed")
	}
}

func TestToXML(t *testing.T) {
	feed := testFeed(t)
	feed.Parse()

	xml, err := feed.ToXML()
	if err != nil || xml == "" {
		t.Errorf("Unable to write feed to XML")
	}
}

func TestFlow(t *testing.T) {
	feed := testFeed(t)
	err := feed.Parse()
	if err != nil {
		t.Errorf("Unable to parse feed")
	}
	err = feed.Filter()

	fmt.Printf("%v\n", feed.FilteredItems)

	if err != nil {
		t.Errorf("Unable to filter feed")
	}

	xml, err := feed.ToXML()
	if err != nil || xml == "" {
		t.Errorf("Unable to write feed to XML")
	}

	fmt.Printf("%s\n", xml)

}
