FROM golang:1.7.4

RUN go get -v github.com/garyburd/redigo/redis
ADD src/ /go/src/
RUN go get -v ./...
CMD "server"
EXPOSE "8080"
