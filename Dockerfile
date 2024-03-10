FROM golang:1.22.1
MAINTAINER Jian-Long Huang <contact@jlhg.org>
EXPOSE 8080
ENV LANG=C.UTF-8
ENV GIN_MODE=release
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN make
CMD ./bin/webserver
