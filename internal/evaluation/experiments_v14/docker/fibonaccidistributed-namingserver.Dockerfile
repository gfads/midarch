# Building app
FROM golang:1.20-bullseye as build
WORKDIR /go

ENV GOPATH /go
ENV GOROOT /usr/local/go
ENV GMIDARCHDIR /go
ENV GMIDARCH /go
ENV FDR4 /usr/local/fdr4/bin
RUN go env -w GO111MODULE=off

# COPY ./fdr4 /usr/local/fdr4
# #RUN sh -c 'echo "deb http://www.cs.ox.ac.uk/projects/fdr/downloads/debian/ fdr release\n" > /etc/apt/sources.list.d/fdr.list'
# #RUN wget -qO - http://www.cs.ox.ac.uk/projects/fdr/downloads/linux_deploy.key | apt-key add -
# RUN apt update
# #RUN apt-cache showpkg fdr
# RUN apt install libfontconfig1 libfreetype6 libice6 libsm6 libx11-6 libxau6 libxdmcp6 libxext6 libxrender1 -y
# #COPY ./libpng12.so.0 /usr/lib/fdr4/libpng12.so.0
# #RUN echo "/usr/lib/fdr4" > /etc/ld.so.conf.d/libpng12.conf
# #RUN apt install fdr -y

#RUN go get gopkg.in/check.v1
#RUN go get github.com/kr/pretty
#RUN go get github.com/kr/text
#RUN go get github.com/vmihailenco/msgpack
#RUN go get github.com/vmihailenco/tagparser
#RUN go get github.com/lucas-clemente/quic-go

#RUN go get -u google.golang.org/grpc
#RUN go get -u github.com/golang/protobuf/protoc-gen-go
#RUN go get -u github.com/streadway/amqp

COPY ./pkg ./pkg

RUN go build $GMIDARCHDIR/examples/fibonaccidistributed/naming/naming.go

# add app
#COPY docker-entrypoint.sh /usr/local/bin
#RUN chmod +x /usr/local/bin/docker-entrypoint.sh
EXPOSE 1313
#ENTRYPOINT ["docker-entrypoint.sh"]

CMD ["/go/naming"]
