FROM golang:1.11

LABEL maintainer="Alexander Shinov <alexandershinov@gmail.com>"

WORKDIR $GOPATH/src/github.com/alexandershinov/imageEater

COPY "./*.go" ./
COPY "./*.toml" /opt/imageEater/
COPY "./saver" ./saver
RUN go get -d -v ./...
RUN go get github.com/stretchr/testify/assert
RUN go test ./...
RUN go install -v ./...
