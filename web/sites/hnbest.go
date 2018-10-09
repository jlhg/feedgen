package sites

import (
    "net/http"
    "io/ioutil"
    "log"
    "fmt"
    "time"
    "regexp"
    "strconv"
    "strings"
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
        itemBeforeTime := string(m[5])
        itemBeforeTimeUnit := string(m[6])
        itemCommentPath := string(m[7])
        itemCommentCount := string(m[8])
        itemCommentLink := fmt.Sprintf("https://news.ycombinator.com/%s", itemCommentPath)
        itemDescription := fmt.Sprintf("%s points. <a href=\"%s\" >%s comments</a>", itemPoint, itemCommentLink, itemCommentCount)
        created := now

        if strings.Contains(itemBeforeTimeUnit, "day") {
            day, _ := strconv.Atoi(itemBeforeTime)
            created = now.AddDate(0, 0, -day)
        } else if strings.Contains(itemBeforeTimeUnit, "hour") {
            duration, _ := time.ParseDuration(fmt.Sprintf("-%sh", itemBeforeTime))
            created = now.Add(duration)
        } else if strings.Contains(itemBeforeTimeUnit, "minute") {
            duration, _ := time.ParseDuration(fmt.Sprintf("-%sm", itemBeforeTime))
            created = now.Add(duration)
        }

        feed.Add(&feeds.Item{
            Id: itemCommentLink,
            Title: string(itemTitle),
            Link: &feeds.Link{Href: itemLink},
            Description: itemDescription,
            Author: &feeds.Author{Name: itemAuthor},
            Created: created,
        })
    }

    feedText, err = feed.ToAtom()
    if err != nil {
        log.Fatal(err)
    }

    return
}
