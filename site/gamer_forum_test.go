package site

import (
	"net/http"
	"net/url"
	"testing"
)

// TestGamerForumParser_GetFeed_Pokemon tests feed generation for 精靈寶可夢哈拉區
func TestGamerForumParser_GetFeed_Pokemon(t *testing.T) {
	t.Parallel()
	parser := GamerForumParser{}

	// Test with bsn=1647 (精靈寶可夢)
	query := url.Values{}
	query.Set("bsn", "1647")

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

// TestGamerForumParser_GetFeed_MonsterHunter tests feed generation for 魔物獵人 (20推以上)
func TestGamerForumParser_GetFeed_MonsterHunter(t *testing.T) {
	t.Parallel()
	parser := GamerForumParser{}

	// Test with bsn=5786&gp=20 (魔物獵人 20推以上)
	query := url.Values{}
	query.Set("bsn", "5786")
	query.Set("gp", "20")

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

// TestGamerForumParser_ArticleLinks tests HTTP response for article links
func TestGamerForumParser_ArticleLinks(t *testing.T) {
	t.Parallel()
	parser := GamerForumParser{}
	query := url.Values{}
	query.Set("bsn", "1647")

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

// TestGamerForumParser_MissingParameter tests error handling for missing bsn parameter
func TestGamerForumParser_MissingParameter(t *testing.T) {
	t.Parallel()
	parser := GamerForumParser{}
	query := url.Values{}

	_, err := parser.GetFeed(query)
	if err == nil {
		t.Error("Expected error for missing bsn parameter, got nil")
	}
}

// TestGamerForumParser_InvalidParameter tests error handling for invalid gp parameter
func TestGamerForumParser_InvalidParameter(t *testing.T) {
	t.Parallel()
	parser := GamerForumParser{}
	query := url.Values{}
	query.Set("bsn", "1647")
	query.Set("gp", "999") // Invalid gp value

	_, err := parser.GetFeed(query)
	if err == nil {
		t.Error("Expected error for invalid gp parameter, got nil")
	}
}
