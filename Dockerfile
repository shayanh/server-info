FROM golang:1.6

ADD . /go/src/github.com/shayanh/server-info
RUN go install github.com/shayanh/server-info
ENTRYPOINT /go/bin/server-info

EXPOSE 80