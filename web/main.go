package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "./sites"
)

// Context ...
type Context struct {
    *gin.Context
}

func feedContentTypeMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Content-Type", "application/atom+xml; charset=utf-8")
        c.Next()
    }
}

func setupRouter() *gin.Engine {
	r := gin.Default()
    r.Use(feedContentTypeMiddleware())

    r.GET("/hnbest", func(context *gin.Context) {
        c := &Context{context}
        feedText, err := sites.HNBestFeed()
        if err != nil {
            c.String(http.StatusServiceUnavailable, err.Error())
            return
        }
        c.String(http.StatusOK, feedText)
    })

    r.GET("/ptt/:boardName", func(context *gin.Context) {
        c := &Context{context}
        boardName := c.Param("boardName")
        query := c.Query("q")
        args := &sites.PttArgument{BoardName: boardName, Query: query}
        feedText, err := sites.PttFeed(args)
        if err != nil {
            c.String(http.StatusServiceUnavailable, err.Error())
            return
        }
        c.String(http.StatusOK, feedText)
    })

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run()
}
