package site

import (
    "net/http"
    "io/ioutil"
    "log"
    "fmt"
    "time"
    "regexp"

    "github.com/gin-gonic/gin"
    "github.com/gorilla/feeds"

    "github.com/jlhg/feedgen"
)

// PttRouter is a route handler for https://www.ptt.cc/index.html.
func PttRouter(c *gin.Context) {
    boardName := c.Param("boardName")
    query := c.Query("q")
    feedText, err := getPttFeedText(boardName, query)
    if err != nil {
        log.Println(err)
        c.String(http.StatusServiceUnavailable, err.Error())
        return
    }

    c.Header("Content-Type", "application/atom+xml; charset=utf-8")
    c.String(http.StatusOK, feedText)

    return
}

func getPttFeedText(boardName string, query string) (feedText string, err error) {
    now := time.Now()
    title := fmt.Sprintf("批踢踢實業坊 %s 板", boardName)
    var url string
    if query == "" {
        url = "https://www.ptt.cc/bbs/" + boardName + "/index.html"
    } else {
        url = "https://www.ptt.cc/bbs/" + boardName + "/search?q=" + query
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

    re = regexp.MustCompile(fmt.Sprintf(`<a href="(/bbs/%s/M\..+?\.html)">`, boardName))
    matchGroup := re.FindAllSubmatch(match, -1)
    feedItemsCount := len(matchGroup)
    if feedItemsCount == 0 {
        err = &feedgen.ArticleLinkFetchError{url}
        return
    }

    var feedItems []*feeds.Item
    for _, m := range matchGroup {
        url := "https://www.ptt.cc" + string(m[1])

        var feedItem *feeds.Item
        feedItem, err = getPttArticleFeedItem(url)
        if err != nil {
            return
        }

        feedItems = append(feedItems, feedItem)
    }

    for _, feedItem := range feedItems {
        feed.Add(feedItem)
    }

    feedText, err = feed.ToAtom()
    if err != nil {
        log.Fatal(err)
    }

    return
}


func getPttArticleFeedItem(url string) (feedItem *feeds.Item, err error) {
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
        if match == nil {
            err = &feedgen.ArticleContentFetchError{url}
            return
        }
        description = "<pre>" + string(match[1]) + "</pre>"
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
        content = regexp.MustCompile(`(?s)<div.+?>(.+?)</div>`).ReplaceAllString(content, "$1")
        content = regexp.MustCompile(`(?s)<span.+?>(.+?)</span>`).ReplaceAllString(content, "$1")
        content = regexp.MustCompile(`(?s)<a.+?>(.+?)</a>`).ReplaceAllString(content, "$1")
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

    feedItem = &feeds.Item{
        Id: url,
        Title: title,
        Link: &feeds.Link{Href: url},
        Description: description,
        Author: &feeds.Author{Name: author},
        Created: created,
    }

    return
}
