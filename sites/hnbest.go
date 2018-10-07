package sites

import (
    "net/http"
    "io/ioutil"
    "log"
    "fmt"
    "time"
    "regexp"
    "strconv"
    "github.com/gorilla/feeds"
)

// HNBestFeed ...
func HNBestFeed() (feedText string, err error) {
    now := time.Now()
    title := "Top Links | Hacker News"
    url := "https://news.ycombinator.com/best"
    feed := feeds.Feed{
        Title: title,
        Link: &feeds.Link{Href: url},
        Description: "",
        Author: nil,
        Created: now,
    }

    resp, err := http.Get(url)
    if err != nil {
        return
    }

    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return
    }

    re := regexp.MustCompile(`(?s)<td class="title"><a href="(.+?)" class="storylink">(.+?)</a>.+?<span class="score" id=".+?">(\d+?) points</span>.+?by <a href=".+?" class="hnuser">(.+?)</a>.+?(\d+?) (days?|hours?|minutes?) ago.+?<a href="(.+?)">(\d+?)&nbsp;comments</a>`)
    matchGroup := re.FindAllSubmatch(body, -1)
    for _, m := range matchGroup {
        itemLink := string(m[1])
        itemTitle := string(m[2])
        itemPoint := string(m[3])
        itemAuthor := string(m[4])
        itemBeforeDay, _ := strconv.Atoi(string(m[5]))
        itemCommentPath := string(m[7])
        itemCommentCount := string(m[8])
        itemDescription := fmt.Sprintf("%s points. <a href=\"https://news.ycombinator.com/%s\" >%s comments</a>", itemPoint, itemCommentPath, itemCommentCount)

        feed.Add(&feeds.Item{
            Title: string(itemTitle),
            Link: &feeds.Link{Href: itemLink},
            Description: itemDescription,
            Author: &feeds.Author{Name: itemAuthor},
            Created: now.AddDate(0, 0, -itemBeforeDay),
        })
    }

    feedText, err = feed.ToAtom()
    if err != nil {
        log.Fatal(err)
    }

    return
}
