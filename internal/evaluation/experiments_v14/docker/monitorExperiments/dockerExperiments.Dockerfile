# Building app
FROM golang:1.15.9-buster as build
WORKDIR /go

ENV GOPATH /go
ENV GOROOT /usr/local/go

#RUN go get github.com/moby/moby

COPY ./evaluation/experiments/docker/monitorExperiments ./pkg
COPY ./pkg/github.com ./pkg/github.com

RUN go build $GOPATH/src/dockerExperiments.go

RUN sysctl -w net.core.rmem_max=7500000
RUN sysctl -w net.core.wmem_max=7500000

CMD ["/go/dockerExperiments"]
