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

func (c *Context) setHeader() {
    c.Header("Content-Type", "application/atom+xml; charset=utf-8")
}

func setupRouter() *gin.Engine {
	r := gin.Default()

    r.GET("/hnbest", func(context *gin.Context) {
        c := &Context{context}
        c.setHeader()
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
        c.setHeader()
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
