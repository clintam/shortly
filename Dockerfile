FROM golang:1.8

RUN go get -v \
           github.com/garyburd/redigo/redis \
           gopkg.in/mgo.v2

ADD src/ /go/src/
RUN go get -v ./src/server ./src/test-suite

CMD "server"
EXPOSE "8080"
