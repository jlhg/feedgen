FROM golang:1.11
MAINTAINER Jian-Long Huang
EXPOSE 8080
RUN mkdir /web
COPY ./web /web
WORKDIR /web
RUN go get github.com/gin-gonic/gin github.com/gorilla/feeds
CMD go run main.go
