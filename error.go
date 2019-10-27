package feedgen

import (
    "fmt"
)

// ArticleLinkFetchError shows that the article links can't be fetched from the source URL.
type ArticleLinkFetchError struct {
    SourceURL string
}

func (e *ArticleLinkFetchError) Error() string {
    return fmt.Sprintf("article links can't be fetched from the source %s", e.SourceURL)
}

// ArticleContentFetchError shows that the article content can't be fetched from the source URL.
type ArticleContentFetchError struct {
    SourceURL string
}

func (e *ArticleContentFetchError) Error() string {
    return fmt.Sprintf("article content can't be fetched from the source %s", e.SourceURL)
}
