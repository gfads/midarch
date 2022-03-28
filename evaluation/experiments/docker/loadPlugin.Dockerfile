# Building app
FROM golang:1.15.9-buster as build
WORKDIR /go

ENV GOPATH /go
ENV GOROOT /usr/local/go
ENV GMIDARCH /go
ENV GMIDARCHDIR /go
ENV FDR4 /usr/local/fdr4/bin

#RUN go get google.golang.org/grpc

COPY ./src ./src
COPY ./evaluation/experiments ./evaluation/experiments

RUN go build $GMIDARCHDIR/src/adaptive/loadPlugin/loadPlugin.go

CMD ["/go/loadPlugin"]
