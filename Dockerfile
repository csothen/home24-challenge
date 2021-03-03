FROM golang:1.15-alpine AS build

RUN apk --no-cache add gcc g++ make git
WORKDIR /go/src/github.com/csothen/htmlparser
COPY . .
RUN go get ./...
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/server ./cmd/htmlparser

EXPOSE 9090
ENTRYPOINT /go/src/github.com/csothen/htmlparser/bin/server