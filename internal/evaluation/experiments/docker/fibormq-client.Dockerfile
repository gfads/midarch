# Building app
FROM golang:1.15.9-buster as build
WORKDIR /go

ENV GOPATH /go
ENV GOROOT /usr/local/go
ENV GMIDARCHDIR /go
ENV FDR4 /usr/local/fdr4/bin

RUN go get github.com/streadway/amqp

COPY ./pkg ./pkg
COPY ./evaluation/experiments ./evaluation/experiments

RUN go build $GMIDARCHDIR/evaluation/experiments/fiboApps/fibo_rmq/client/client.go

CMD ["/go/client"]