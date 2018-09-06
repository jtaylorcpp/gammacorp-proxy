FROM golang:alpine

RUN apk add git -y

ADD . $GOPATH/src/github.com/jtaylorcpp/gammacorp-proxy/
WORKDIR $GOPATH/src/github.com/jtaylorcpp/gammacorp-proxy/

RUN go get ...

WORKDIR $GOPATH/src/github.com/jtaylorcpp/gammacorp-proxy/cmd/simple

RUN go build -o /usr/bin/simpleproxy

CMD /usr/bin/simpleproxy
