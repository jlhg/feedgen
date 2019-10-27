package main

import (
    "io"
    "log"
    "os"
    "path"

    "github.com/gin-gonic/gin"

    "github.com/jlhg/feedgen/site"
)

func cache() gin.HandlerFunc {
    return func(c *gin.Context) {
        // TODO
    }
}

func setRouter() *gin.Engine {
	r := gin.Default()
    r.Use(cache())

    r.GET("/hackernews/:category", site.HackerNewsRouter)
    r.GET("/ptt/:boardName", site.PttRouter)

	return r
}

func setLogger() {
    gin.DisableConsoleColor()

    var fileName string
    if os.Getenv("GIN_MODE") == "release" {
        fileName = "production.log"
    } else {
        fileName = "development.log"
    }

    dir := "log"
    filePath := path.Join("log", fileName)

    if _, err := os.Stat(dir); os.IsNotExist(err) {
        if err := os.Mkdir(dir, os.ModePerm); err != nil {
            panic(err)
        }
    }

    f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        panic(err)
    }

    gin.DefaultWriter = io.MultiWriter(f)
    log.SetOutput(gin.DefaultWriter)
}

func main() {
    setLogger()
	r := setRouter()
	r.Run()
}
