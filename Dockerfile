FROM alpine:3.6

RUN apk add --no-cache  ca-certificates

ADD bin/linux/amd64/parser /parser
EXPOSE 8083
ENTRYPOINT ["/parser"]
