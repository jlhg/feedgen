package feedgen

import (
	"net/url"

	"github.com/gorilla/feeds"
)

// QueryValues is an alias of url.Values. It is typically used for query parameters and form values.
type QueryValues = url.Values

// Parser is a interface that defines a method to generates feed from query parameters.
type Parser interface {
	// GetFeed returns feed generated from the source site.
	GetFeed(query QueryValues) (feed *feeds.Feed, err error)
}
