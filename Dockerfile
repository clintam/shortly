FROM golang:1.7.4

ADD src/ /go/src/
RUN go get -v ./...
CMD "server"
EXPOSE "8080"
