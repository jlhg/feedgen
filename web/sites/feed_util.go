package sites

import (
    "sort"
    "github.com/gorilla/feeds"
)

// SortFeedItemsLatestFirst ...
func SortFeedItemsLatestFirst(members []*feeds.Item) {
    sort.Slice(members, func(i, j int) bool { return members[i].Created.After(members[j].Created) })
}
