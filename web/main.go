package main

import (
    "bytes"
    "net/http"
    "io"
    "log"
    "os"
    "path"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis/v7"

    "github.com/jlhg/feedgen/site"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
    w.body.WriteString(s)
    return w.ResponseWriter.WriteString(s)
}

func cache(client *redis.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
        c.Writer = blw

        reqURL := c.Request.URL.String()
        feedText, err := client.Get(reqURL).Result()
        if err == nil {
            c.Header("Content-Type", "application/atom+xml; charset=utf-8")
            c.String(http.StatusOK, feedText)
            c.Abort()
            return
        }

        c.Next()

        if c.Writer.Status() == http.StatusOK {
            err := client.Set(reqURL, blw.body.String(), time.Minute).Err()
            if err != nil {
                panic(err)
            }
        }
    }
}

func setRouter() *gin.Engine {
	r := gin.Default()
    redisHost := os.Getenv("FG_REDIS_HOST")
    redisPassword := os.Getenv("FG_REDIS_PASSWORD")
    redisDB := 0
    if redisDBString := os.Getenv("FG_REDIS_DB"); redisDBString != "" {
        var err error
        redisDB, err = strconv.Atoi(redisDBString)
        if err != nil {
            panic(err)
        }
    }
    if redisHost != "" {
        client := redis.NewClient(&redis.Options{Addr: redisHost, Password: redisPassword, DB: redisDB})
        r.Use(cache(client))
    }

    r.GET("/", func(c *gin.Context) {
        c.Redirect(http.StatusTemporaryRedirect, "https://github.com/jlhg/feedgen")
    })
    r.GET("/hackernews/:category", site.HackerNewsRouter)
    r.GET("/ptt/:boardName", site.PttRouter)
    r.GET("/gamer_forum/:bsn", site.GamerForumRouter)

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
