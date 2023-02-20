package site

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gorilla/feeds"
	"golang.org/x/exp/slices"

	"github.com/jlhg/feedgen"
)

type UdnGameData struct {
	Articles map[string]UdnGameArticles `json:"articles"`
}

type UdnGameArticles struct {
	ArtTitle      string `json:"art_title"`
	ArtAuthorName string `json:"art_author_name"`
	Summary       string `json:"summary"`
	Link          string `json:"link"`
	ArtTime       int64  `json:"art_time"`
}

// UdnGameParser is a parser for 遊戲角落 (https://game.udn.com/rank/newest/2003).
type UdnGameParser struct{}

// GetFeed returns generated feed with the given query parameters.
func (parser UdnGameParser) GetFeed(query feedgen.QueryValues) (feed *feeds.Feed, err error) {
	now := time.Now()

	var sourceLink, rawLink string
	var title, subTitle string

	section := query.Get("section")
	switch section {
	case "rank":
		by := query.Get("by")
		switch by {
		case "newest":
			sourceLink = "https://game.udn.com/rank/newest/2003"
			rawLink = "https://game.udn.com/rank/ajax_newest/2003/0/1"
			subTitle = "最新文章"
		case "pv":
			sourceLink = "https://game.udn.com/rank/pv/2003"
			rawLink = "https://game.udn.com/rank/ajax_pv/2003/0/2"
			subTitle = "最多瀏覽"
		default:
			err = &feedgen.ParameterValueInvalidError{"by"}
		}
	default:
		err = &feedgen.ParameterValueInvalidError{"section"}
	}

	title = fmt.Sprintf("%s | 遊戲角落", subTitle)

	feed = &feeds.Feed{
		Title:       title,
		Link:        &feeds.Link{Href: sourceLink},
		Description: "",
		Author:      nil,
		Created:     now,
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", rawLink, nil)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	data := UdnGameData{}
	err = sonic.Unmarshal(body, &data)
	if err != nil {
		return
	}

	keys := []int{}

	for k := range data.Articles {
		var i int
		i, err = strconv.Atoi(k)
		if err != nil {
			return
		}

		keys = append(keys, i)
	}

	slices.Sort(keys)

	for _, i := range keys {
		var created time.Time

		a := fmt.Sprintf("%d", i)
		article := data.Articles[a]

		switch feedgen.CountDigits(article.ArtTime) {
		case 10:
			created = time.Unix(article.ArtTime, 0)
		case 13:
			created = time.UnixMilli(article.ArtTime)
		default:
			err = &feedgen.ItemFetchError{rawLink}
			return
		}

		feedItem := &feeds.Item{
			Id:          article.Link,
			Title:       article.ArtTitle,
			Link:        &feeds.Link{Href: article.Link},
			Description: article.Summary,
			Author:      &feeds.Author{Name: article.ArtAuthorName},
			Created:     created,
		}

		feed.Add(feedItem)
	}

	return
}
