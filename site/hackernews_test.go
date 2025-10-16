package site

import (
	"net/url"
	"testing"
)

// TestHackernewsParser_GetFeed_Best tests feed generation for Hacker News best category
func TestHackernewsParser_GetFeed_Best(t *testing.T) {
	t.Parallel()
	parser := HackernewsParser{}

	// Test with category=best
	query := url.Values{}
	query.Set("category", "best")

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

// TestHackernewsParser_ArticleLinks tests URL format for article links
func TestHackernewsParser_ArticleLinks(t *testing.T) {
	t.Parallel()
	parser := HackernewsParser{}
	query := url.Values{}
	query.Set("category", "best")

	feed, err := parser.GetFeed(query)
	if err != nil {
		t.Fatalf("Failed to get feed: %v", err)
	}

	if len(feed.Items) == 0 {
		t.Skip("No items to test")
	}

	for i, item := range feed.Items {
		if item.Link == nil || item.Link.Href == "" {
			t.Errorf("Item %d has no link", i)
			continue
		}

		// Hacker News links point to external sites, only validate URL format
		_, err := url.ParseRequestURI(item.Link.Href)
		if err != nil {
			t.Errorf("Item %d has invalid URL %s: %v", i, item.Link.Href, err)
		} else {
			t.Logf("Item %d URL is valid: %s", i, item.Link.Href)
		}
	}
}

// TestHackernewsParser_MissingParameter tests error handling for missing category parameter
func TestHackernewsParser_MissingParameter(t *testing.T) {
	t.Parallel()
	parser := HackernewsParser{}
	query := url.Values{}

	_, err := parser.GetFeed(query)
	if err == nil {
		t.Error("Expected error for missing category parameter, got nil")
	}
}

// TestHackernewsParser_InvalidParameter tests error handling for invalid category parameter
func TestHackernewsParser_InvalidParameter(t *testing.T) {
	t.Parallel()
	parser := HackernewsParser{}
	query := url.Values{}
	query.Set("category", "invalid")

	_, err := parser.GetFeed(query)
	if err == nil {
		t.Error("Expected error for invalid category parameter, got nil")
	}
}
