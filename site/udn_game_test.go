package site

import (
	"net/http"
	"net/url"
	"testing"
)

// TestUdnGameParser_GetFeed_Newest tests feed generation for 遊戲角落 newest articles
func TestUdnGameParser_GetFeed_Newest(t *testing.T) {
	t.Parallel()
	parser := UdnGameParser{}

	// Test with section=rank&by=newest
	query := url.Values{}
	query.Set("section", "rank")
	query.Set("by", "newest")

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

// TestUdnGameParser_GetFeed_MostViewed tests feed generation for 遊戲角落 most viewed articles
func TestUdnGameParser_GetFeed_MostViewed(t *testing.T) {
	t.Parallel()
	parser := UdnGameParser{}

	// Test with section=rank&by=pv
	query := url.Values{}
	query.Set("section", "rank")
	query.Set("by", "pv")

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

// TestUdnGameParser_ArticleLinks tests HTTP response for article links
func TestUdnGameParser_ArticleLinks(t *testing.T) {
	t.Parallel()
	parser := UdnGameParser{}
	query := url.Values{}
	query.Set("section", "rank")
	query.Set("by", "newest")

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

// TestUdnGameParser_InvalidSection tests error handling for invalid section parameter
func TestUdnGameParser_InvalidSection(t *testing.T) {
	t.Parallel()
	parser := UdnGameParser{}
	query := url.Values{}
	query.Set("section", "invalid")
	query.Set("by", "newest")

	_, err := parser.GetFeed(query)
	if err == nil {
		t.Error("Expected error for invalid section parameter, got nil")
	}
}

// TestUdnGameParser_InvalidBy tests error handling for invalid by parameter
func TestUdnGameParser_InvalidBy(t *testing.T) {
	t.Parallel()
	parser := UdnGameParser{}
	query := url.Values{}
	query.Set("section", "rank")
	query.Set("by", "invalid")

	_, err := parser.GetFeed(query)
	if err == nil {
		t.Error("Expected error for invalid by parameter, got nil")
	}
}
