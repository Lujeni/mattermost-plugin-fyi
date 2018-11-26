FROM golang:latest as build-stage

ADD https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

WORKDIR /go/src/github.com/lujeni/fyi
ADD . .
RUN make deps & make build
EXPOSE 8888
CMD ["/go/src/github.com/lujeni/fyi/fyi"]
