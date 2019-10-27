package feedgen

import (
    "sort"

    "github.com/gorilla/feeds"
)

// SortFeedItemsLatestFirst sorts feed items by Created in descending order.
func SortFeedItemsLatestFirst(members []*feeds.Item) {
    sort.Slice(members, func(i, j int) bool { return members[i].Created.After(members[j].Created) })
}
