FROM golang:1.12
MAINTAINER Arthur Nguyen <nguyendhoan9@gmail.com>
ADD main.go /go/src/main.go
EXPOSE 8080
WORKDIR /go/src
ENTRYPOINT [ "go", "run", "main.go"]
