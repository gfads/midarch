# Building app
FROM golang:1.15.9-buster as build
WORKDIR /go

ENV GOPATH /go
ENV GOROOT /usr/local/go
ENV GMIDARCHDIR /go
ENV FDR4 /usr/local/fdr4/bin

COPY ./pkg ./pkg
COPY ./evaluation/experiments ./evaluation/experiments

RUN go build $GMIDARCHDIR/evaluation/experiments/fiboApps/fiboRPC/server/server.go

EXPOSE 2030
CMD ["/go/server"]
