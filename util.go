package feedgen

import (
	"sort"

	"github.com/gorilla/feeds"
)

// SortFeedItemsLatestFirst sorts feed items by Created in descending order.
func SortFeedItemsLatestFirst(members []*feeds.Item) {
	sort.Slice(members, func(i, j int) bool { return members[i].Created.After(members[j].Created) })
}

func CountDigits(i int64) (count int) {
	for i != 0 {
		i /= 10
		count = count + 1
	}

	return count
}
