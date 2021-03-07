FROM golang:1.15.3

RUN mkdir -p /go/src/application 

WORKDIR /go/src/application 

ADD . .

RUN go get ./... && go mod vendor && go mod verify

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

CMD ["./run.sh"]