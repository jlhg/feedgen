package site

import (
    "net/http"
    "log"
    "fmt"
    "time"
    "regexp"

    "github.com/gin-gonic/gin"
    "github.com/gorilla/feeds"
    "github.com/PuerkitoBio/goquery"
)

// GamerForumRouter is a route handler for https://forum.gamer.com.tw/.
func GamerForumRouter(c *gin.Context) {
    bsn := c.Param("bsn")
    if matched, _ := regexp.MatchString(`^\d+$`, bsn); !matched {
        c.String(http.StatusBadRequest, "bsn has invalid value")
        return
    }

    gp := c.Query("gp")
    if gp != "" {
        if matched, _ := regexp.MatchString(`^5|20|50|100|200$`, gp); !matched {
            c.String(http.StatusBadRequest, "gp has invalid value")
            return
        }
    }

    feedText, err := getGamerForumFeedText(bsn, gp)
    if err != nil {
        log.Println(err)
        c.String(http.StatusServiceUnavailable, err.Error())
        return
    }

    c.Header("Content-Type", "application/atom+xml; charset=utf-8")
    c.String(http.StatusOK, feedText)

    return
}

func getGamerForumFeedText(bsn string, gp string) (feedText string, err error) {
    now := time.Now()
    url := fmt.Sprintf("https://forum.gamer.com.tw/B.php?bsn=%s", bsn)
    if gp != "" {
        url = fmt.Sprintf("%s&qt=4&q=%s", url, gp)
    }

    client := &http.Client{}
    cookie := http.Cookie{Name: "ckForumListOrder", Value: "post"}
    req, err := http.NewRequest("GET", url, nil)
    req.AddCookie(&cookie)
    resp, err := client.Do(req)
    if err != nil {
        return
    }
    defer resp.Body.Close()

    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return
    }

    title := doc.Find("head > title").Text()
    description, _ := doc.Find(`meta[name="Description"]`).Attr("content")

    feed := feeds.Feed{
        Title: title,
        Link: &feeds.Link{Href: url},
        Description: description,
        Author: nil,
        Created: now,
    }

    doc.Find("[class=\"b-list__row b-list-item b-imglist-item\"]").Each(func(i int, s *goquery.Selection) {
        itemID, _ := s.Find(".b-list__main__title").Attr("href")
        itemTitle := s.Find(".b-list__main__title").Text()
        itemContent := s.Find(".b-list__brief").Text()
        itemAuthor := s.Find(".b-list__count__user > a").Text()
        itemLink := fmt.Sprintf("https://forum.gamer.com.tw/%s", itemID)
        if itemContent == "" {
            itemContent = "(沒有內容)"
        }

        doc, err := goquery.NewDocument(itemLink)
        if err != nil {
            return
        }

        dt, _ := doc.Find(".edittime.tippy-post-info").Attr("data-mtime")
        layout := "2006-01-02 15:04:05"
        itemCreated, err := time.Parse(layout, dt)
        if err != nil {
            return
        }

        feed.Add(&feeds.Item{
            Id: itemID,
            Title: itemTitle,
            Link: &feeds.Link{Href: itemLink},
            Description: itemContent,
            Author: &feeds.Author{Name: itemAuthor},
            Created: itemCreated,
        })
    })

    feedText, err = feed.ToAtom()
    if err != nil {
        log.Fatal(err)
    }

    return
}
