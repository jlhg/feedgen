package sites

import (
    "net/http"
    "io/ioutil"
    "log"
    "fmt"
    "time"
    "regexp"
    "github.com/gorilla/feeds"
)

// PttArgument ...
type PttArgument struct {
    BoardName string
    Query string
}

// PttFeed ...
func PttFeed(args *PttArgument) (feedText string, err error) {
    now := time.Now()
    title := fmt.Sprintf("批踢踢實業坊 %s 板", args.BoardName)
    var url string
    if args.Query == "" {
        url = "https://www.ptt.cc/bbs/" + args.BoardName + "/index.html"
    } else {
        url = "https://www.ptt.cc/bbs/" + args.BoardName + "/search?q=" + args.Query
    }

    feed := feeds.Feed{
        Title: title,
        Link: &feeds.Link{Href: url},
        Description: "",
        Author: nil,
        Created: now,
    }

    client := &http.Client{}
    cookie := http.Cookie{Name: "over18", Value: "1"}
    req, err := http.NewRequest("GET", url, nil)
    req.AddCookie(&cookie)
    resp, err := client.Do(req)
    if err != nil {
        return
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return
    }

    // 排除置底文章
    re := regexp.MustCompile(`(?s)(.+?)<div class="r-list-sep"></div>`)
    match := re.Find(body)
    if match == nil {
        match = body
    }

    re = regexp.MustCompile(fmt.Sprintf(`<a href="(/bbs/%s/M\..+?\.html)">`, args.BoardName))
    matchGroup := re.FindAllSubmatch(match, -1)
    feedItemsCount := len(matchGroup)
    ch := make(chan *feeds.Item)

    for _, m := range matchGroup {
        url := "https://www.ptt.cc" + string(m[1])
        go fetchFeedItem(url, ch)
    }

    var feedItems []*feeds.Item
    count := 0
    for feedItem := range ch {
        count++
        feedItems = append(feedItems, feedItem)
        if count >= feedItemsCount {
            break
        }
    }
    close(ch)
    SortFeedItemsLatestFirst(feedItems)

    for _, feedItem := range feedItems {
        feed.Add(feedItem)
    }

    feedText, err = feed.ToAtom()
    if err != nil {
        log.Fatal(err)
    }

    return
}


func fetchFeedItem(url string, ch chan *feeds.Item) {
    re := regexp.MustCompile(`(?s)<div id="main-content" class="bbs-screen bbs-content"><div class="article-metaline"><span class="article-meta-tag">作者</span><span class="article-meta-value">(.+?)</span></div>(<div class="article-metaline-right"><span class="article-meta-tag">看板</span><span class="article-meta-value">(.+?)</span></div>)?<div class="article-metaline"><span class="article-meta-tag">標題</span><span class="article-meta-value">(.+?)</span></div>(<div class="article-metaline"><span class="article-meta-tag">時間</span><span class="article-meta-value">(.+?)</span></div>)?(.+?)<span class="f2">※ (發信站|編輯)`)
    client := &http.Client{}
    cookie := http.Cookie{Name: "over18", Value: "1"}
    req, err := http.NewRequest("GET", url, nil)
    req.AddCookie(&cookie)
    resp, err := client.Do(req)
    if err != nil {
        return
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return
    }
    match := re.FindSubmatch(body)
    author := string(match[1])
    title := string(match[4])
    const timeForm = "Mon Jan 02 15:04:05 2006"
    created, _ := time.Parse(timeForm, string(match[6]))
    description := string(match[7])

    ch <- &feeds.Item{
        Id: url,
        Title: string(title),
        Link: &feeds.Link{Href: url},
        Description: description,
        Author: &feeds.Author{Name: author},
        Created: created,
    }
}
