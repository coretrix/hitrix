FROM golang:1.17-alpine

RUN mkdir -p /go/src/github.com/coretrix/hitrix

WORKDIR /go/src/github.com/coretrix/hitrix

ADD services/docker-entrypoint.sh /usr/bin/docker-entrypoint
RUN chmod +x /usr/bin/docker-entrypoint

ENTRYPOINT ["docker-entrypoint"]
