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
        go fetchPTTFeedItem(url, ch)
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


func fetchPTTFeedItem(url string, ch chan *feeds.Item) {
    re := regexp.MustCompile(`(?s)<div id="main-content" class="bbs-screen bbs-content"><div class="article-metaline"><span class="article-meta-tag">作者</span><span class="article-meta-value">(.+?)</span></div>(<div class="article-metaline-right"><span class="article-meta-tag">看板</span><span class="article-meta-value">(.+?)</span></div>)?<div class="article-metaline"><span class="article-meta-tag">標題</span><span class="article-meta-value">(.+?)</span></div>(<div class="article-metaline"><span class="article-meta-tag">時間</span><span class="article-meta-value">(.+?)</span></div>)?(.+?)<span class="f2">※ (發信站|編輯)`)
    re2 := regexp.MustCompile(`(?s)class="bbs-screen bbs-content">(.+?)<span class="f2">※ (發信站|編輯)`)
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

    author := "null"
    board := "null"
    title := "null"
    created := time.Now()
    description := "null"

    match := re.FindSubmatch(body)
    if match == nil {
        match = re2.FindSubmatch(body)
        if match != nil {
            description = "<pre>" + string(match[1]) + "</pre>"
        }
    } else {
        author = string(match[1])
        board = string(match[3])
        title = string(match[4])
        const timeForm = "Mon Jan 2 15:04:05 2006"
        date := string(match[6])
        created, _ = time.Parse(timeForm, date)
        content := string(match[7])
        content = regexp.MustCompile(`(?s)<div class="richcontent"><blockquote.+?</script></div>`).ReplaceAllString(content, "")
        content = regexp.MustCompile(`(?s)<div class="richcontent"><div class="resize-container"><div class="resize-content"><iframe.+</iframe></div></div></div>`).ReplaceAllString(content, "")
        content = regexp.MustCompile(`(?s)<div class="richcontent"><img src=".+?" alt="" /></div>`).ReplaceAllString(content, "")
        description = "<pre>"
        if board != "" {
            description += "看板：" + board + "\n"
        }
        description += "作者：" + author + "\n" + "標題：" + title + "\n"
        if date != "" {
            description += "時間：" + date + "\n"
        }
        description += content + "</pre>"
    }

    ch <- &feeds.Item{
        Id: url,
        Title: title,
        Link: &feeds.Link{Href: url},
        Description: description,
        Author: &feeds.Author{Name: author},
        Created: created,
    }
}
