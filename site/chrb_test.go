package site

import (
	"net/http"
	"net/url"
	"testing"
)

// TestChrbParser_GetFeed tests feed generation for 大管家房屋網
func TestChrbParser_GetFeed(t *testing.T) {
	t.Parallel()
	parser := ChrbParser{}

	// Test with default query (as shown in README)
	query := url.Values{}

	feed, err := parser.GetFeed(query)
	if err != nil {
		t.Fatalf("Failed to get feed: %v", err)
	}

	if feed == nil {
		t.Fatal("Feed is nil")
	}

	if feed.Title == "" {
		t.Error("Feed title is empty")
	}

	if len(feed.Items) == 0 {
		t.Error("Feed has no items")
	}

	t.Logf("Successfully generated feed with %d items", len(feed.Items))
}

// TestChrbParser_ArticleLinks tests HTTP response for article links
func TestChrbParser_ArticleLinks(t *testing.T) {
	t.Parallel()
	parser := ChrbParser{}
	query := url.Values{}

	feed, err := parser.GetFeed(query)
	if err != nil {
		t.Fatalf("Failed to get feed: %v", err)
	}

	if len(feed.Items) == 0 {
		t.Skip("No items to test")
	}

	// Test only the first article link
	item := feed.Items[0]
	if item.Link == nil || item.Link.Href == "" {
		t.Error("First item has no link")
		return
	}

	client := &http.Client{}
	resp, err := client.Get(item.Link.Href)
	if err != nil {
		t.Errorf("Failed to access article link %s: %v", item.Link.Href, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		t.Errorf("Article link %s returned non-2XX status: %d", item.Link.Href, resp.StatusCode)
	} else {
		t.Logf("Article link %s returned status %d", item.Link.Href, resp.StatusCode)
	}
}
