FROM golang:alpine AS builder
ADD . ${GOPATH}/src/github.com/Haujilo/httpbin-go/
WORKDIR ${GOPATH}/src/github.com/Haujilo/httpbin-go/
RUN apk add --no-cache build-base && go test -v ./... && go build -v *.go

FROM alpine:latest
COPY --from=builder /go/src/github.com/Haujilo/httpbin-go/httpbin /usr/local/bin/
EXPOSE 80/tcp
CMD httpbin 0.0.0.0:80
