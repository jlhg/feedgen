FROM golang:1.13
MAINTAINER Jian-Long Huang <huang@jianlong.org>
EXPOSE 8080
ENV LANG=C.UTF-8
ENV GIN_MODE=release
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN make
CMD ./bin/webserver
