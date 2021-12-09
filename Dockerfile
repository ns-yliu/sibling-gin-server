# syntax=docker/dockerfile:1
FROM golang:1.17-alpine3.14 AS builder

RUN --mount=type=cache,target=/var/cache/apk \
    apk add bash ca-certificates git build-base

WORKDIR /app
COPY . .

ENV GOOS=linux
ENV GOARCH=amd64

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -v -tags netgo -ldflags '-w -extldflags "-static"' ./


FROM alpine:3.14.3 as runtime

RUN apk add --no-cache ca-certificates

COPY --from=builder / /ns/edlp/

WORKDIR /ns/edlp

CMD ["./sibling-gin-server"]